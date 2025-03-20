package utils

import (
	"context"
	"fmt"
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

// ParseVersion parses a version string (e.g., "2.5.3" or "v2.5.3")
func ParseVersion(v string) (Version, error) {
	v = strings.TrimPrefix(v, "v") // Remove "v" prefix
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("invalid version format: %s", v)
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
	return current.Major == target.Major
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
