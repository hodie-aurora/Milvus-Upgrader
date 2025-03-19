package upgrade

import (
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"github.com/hodie-aurora/milvus-upgrader/pkg/utils"
)

// Upgrade performs a Milvus upgrade
func Upgrade(instance, namespace, sourceVersion, targetVersion string, force, skipChecks bool, kubeconfig string) error {
	client, err := k8s.GetClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("Failed to get Kubernetes client: %v", err)
	}
	sourceVer, err := utils.ParseVersion(sourceVersion)
	if err != nil {
		return fmt.Errorf("Failed to parse source version: %v", err)
	}
	targetVer, err := utils.ParseVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("Failed to parse target version: %v", err)
	}
	fmt.Printf("Upgrading %s (%s) from %s to %s\n", instance, namespace, sourceVersion, targetVersion)
	if force {
		fmt.Println("Force upgrade enabled")
	}
	if skipChecks {
		fmt.Println("Pre-upgrade checks skipped")
	}
	if utils.IsMinorUpgrade(sourceVer, targetVer) {
		return MinorUpgrade(client, instance, namespace, targetVersion)
	}
	return MajorUpgrade(client, instance, namespace, targetVersion)
}
