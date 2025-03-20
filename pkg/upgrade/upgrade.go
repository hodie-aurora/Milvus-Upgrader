package upgrade

import (
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"github.com/hodie-aurora/milvus-upgrader/pkg/utils"
)

// Upgrade performs a Milvus upgrade
func Upgrade(instance, namespace, sourceVersion, targetVersion string, force, skipChecks bool, kubeconfig string) error {
	// Get Kubernetes client
	client, err := k8s.GetClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %v", err)
	}

	// If sourceVersion is not provided, detect the current version
	if sourceVersion == "" {
		sourceVersion, err = utils.GetCurrentVersion(client, namespace, instance)
		if err != nil {
			return fmt.Errorf("failed to get current version: %v", err)
		}
		fmt.Printf("Detected current version: %s\n", sourceVersion)
	}

	// Parse source version
	sourceVer, err := utils.ParseVersion(sourceVersion)
	if err != nil {
		return fmt.Errorf("failed to parse source version: %v", err)
	}

	// Parse target version
	targetVer, err := utils.ParseVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("failed to parse target version: %v", err)
	}

	// Print upgrade information
	fmt.Printf("Upgrading %s (%s) from %s to %s\n", instance, namespace, sourceVersion, targetVersion)
	if force {
		fmt.Println("Force upgrade enabled")
	}
	if skipChecks {
		fmt.Println("Pre-upgrade checks skipped")
	}

	// Perform minor or major upgrade based on version comparison
	if utils.IsMinorUpgrade(sourceVer, targetVer) {
		return MinorUpgrade(client, instance, namespace, targetVersion)
	}
	return MajorUpgrade(client, instance, namespace, targetVersion)
}
