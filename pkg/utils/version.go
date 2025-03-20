package utils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Version represents a version number
type Version struct {
	Major int
	Minor int
	Patch int
}

// ParseVersion parses a version string (e.g., "2.5.3", "v2.5.3", "3.5.18-r1")
func ParseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v") // Remove "v" prefix
	// Split on non-numeric parts (e.g., "3.5.18-r1" -> ["3", "5", "18"])
	parts := strings.FieldsFunc(v, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	if len(parts) < 3 {
		return Version{}, fmt.Errorf("invalid version format: %s (expected at least x.y.z)", v)
	}
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("invalid major version: %s", parts[0])
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor version: %s", parts[1])
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid patch version: %s", parts[2])
	}
	return Version{Major: major, Minor: minor, Patch: patch}, nil
}

// IsMinorUpgrade determines if it is a minor version upgrade
func IsMinorUpgrade(current, target Version) bool {
	return current.Major == target.Major && current.Minor == target.Minor
}

// GetCurrentVersion gets the current Milvus version from the cluster
func GetCurrentVersion(client *k8s.ClientSet, namespace, instance string) (string, error) {
	// Define the GroupVersionResource for Milvus CRD
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}

	// Get the Milvus Custom Resource
	obj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), instance, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("error: Milvus instance %s not found in namespace %s", instance, namespace)
		}
		return "", fmt.Errorf("failed to get Milvus CR: %v", err)
	}

	// Extract the image field from spec.components.image
	image, found, err := unstructured.NestedString(obj.UnstructuredContent(), "spec", "components", "image")
	if !found || err != nil {
		return "", fmt.Errorf("failed to get image from Milvus CR: %v", err)
	}

	// Assume image is in the form "milvusdb/milvus:2.5.3", extract the version
	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid image format: %s", image)
	}
	version := parts[1]

	// Trim "v" prefix if present to ensure consistency
	version = strings.TrimPrefix(version, "v")

	return version, nil
}

// CheckDependencies checks if current dependency versions meet the target version requirements
func CheckDependencies(client *k8s.ClientSet, namespace, instance, targetVersion string) error {
	// Parse target version
	targetVer, err := ParseVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("failed to parse target version %s: %v", targetVersion, err)
	}

	// Define required dependency versions based on target Milvus version
	var requiredPulsar, requiredEtcd string
	switch {
	case targetVer.Major == 2 && targetVer.Minor <= 4: // 2.2.x to 2.4.x
		requiredPulsar = "2.8.0"
		requiredEtcd = "3.5.0"
	case targetVer.Major == 2 && targetVer.Minor >= 5: // 2.5.x and above
		requiredPulsar = "3.0.0"
		requiredEtcd = "3.5.16"
	default:
		return fmt.Errorf("unsupported target version: %s", targetVersion)
	}

	// Get current Pulsar version
	pulsarVer, err := getDependencyVersion(client, namespace, instance, "pulsar")
	if err != nil {
		return fmt.Errorf("failed to get Pulsar version: %v", err)
	}
	fmt.Printf("Current Pulsar version: %s, required version: >= %s\n", pulsarVer, requiredPulsar)
	if pulsarVer == "unknown" {
		if !promptUserConfirmation("Pulsar version is unknown. Please confirm if your Pulsar version meets the requirement (>= " + requiredPulsar + "). Continue with upgrade? [yes/no]: ") {
			return fmt.Errorf("user chose to abort upgrade due to unknown Pulsar version")
		}
	} else if !isVersionCompatible(pulsarVer, requiredPulsar) {
		return fmt.Errorf("error: Pulsar version %s does not meet the requirement (>= %s) for target version %s; please upgrade Pulsar", pulsarVer, requiredPulsar, targetVersion)
	}

	// Get current Etcd version
	etcdVer, err := getDependencyVersion(client, namespace, instance, "etcd")
	if err != nil {
		return fmt.Errorf("failed to get Etcd version: %v", err)
	}
	fmt.Printf("Current Etcd version: %s, required version: >= %s\n", etcdVer, requiredEtcd)
	if etcdVer == "unknown" {
		if !promptUserConfirmation("Etcd version is unknown. Please confirm if your Etcd version meets the requirement (>= " + requiredEtcd + "). Continue with upgrade? [yes/no]: ") {
			return fmt.Errorf("user chose to abort upgrade due to unknown Etcd version")
		}
	} else if !isVersionCompatible(etcdVer, requiredEtcd) {
		return fmt.Errorf("error: Etcd version %s does not meet the requirement (>= %s) for target version %s; please upgrade Etcd", etcdVer, requiredEtcd, targetVersion)
	}

	fmt.Println("All dependencies meet the requirements for target version", targetVersion)
	return nil
}

// getDependencyVersion retrieves the version of a specific dependency (e.g., "pulsar" or "etcd")
func getDependencyVersion(client *k8s.ClientSet, namespace, instance, depType string) (string, error) {
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}

	obj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), instance, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return "", fmt.Errorf("error: Milvus instance %s not found in namespace %s", instance, namespace)
		}
		return "", fmt.Errorf("failed to get Milvus CR: %v", err)
	}

	var image string
	var found bool

	if depType == "etcd" {
		// Etcd image is under spec.dependencies.etcd.inCluster.values.image.tag
		path := []string{"spec", "dependencies", "etcd", "inCluster", "values", "image", "tag"}
		image, found, err = unstructured.NestedString(obj.UnstructuredContent(), path...)
	} else if depType == "pulsar" {
		// Pulsar image is under spec.dependencies.pulsar.inCluster.values.images.broker.tag
		path := []string{"spec", "dependencies", "pulsar", "inCluster", "values", "images", "broker", "tag"}
		image, found, err = unstructured.NestedString(obj.UnstructuredContent(), path...)
	}

	if !found || err != nil {
		return "unknown", nil // Return "unknown" if dependency version is not found
	}

	// Return raw version string (e.g., "3.5.18-r1" or "3.0.7")
	return image, nil
}

// isVersionCompatible checks if the current version meets or exceeds the required version
func isVersionCompatible(current, required string) bool {
	if current == "unknown" {
		return true // Handled separately with user prompt
	}

	curVer, err := ParseVersion(current)
	if err != nil {
		fmt.Printf("warning: failed to parse current version %s: %v, treating as incompatible\n", current, err)
		return false // Invalid version treated as incompatible
	}
	reqVer, err := ParseVersion(required)
	if err != nil {
		fmt.Printf("warning: failed to parse required version %s: %v, assuming incompatible\n", required, err)
		return false // Should not happen with predefined required versions
	}

	return curVer.Major > reqVer.Major ||
		(curVer.Major == reqVer.Major && curVer.Minor > reqVer.Minor) ||
		(curVer.Major == reqVer.Major && curVer.Minor == reqVer.Minor && curVer.Patch >= reqVer.Patch)
}

// promptUserConfirmation prompts the user with a yes/no question and returns true if "yes"
func promptUserConfirmation(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "yes" || response == "y"
}
