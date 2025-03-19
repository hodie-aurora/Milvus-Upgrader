package upgrade

import (
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
)

// MajorUpgrade performs a major version upgrade
func MajorUpgrade(client *k8s.ClientSet, instance, namespace, targetVersion string) error {
	fmt.Printf("Major version upgrade to %s not yet implemented.\n", targetVersion)
	return fmt.Errorf("Major version upgrade not yet supported")
}
