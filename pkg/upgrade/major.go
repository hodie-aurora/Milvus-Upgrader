package upgrade

import (
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
)

// MajorUpgrade 执行大版本升级
func MajorUpgrade(client *k8s.ClientSet, instance, namespace, targetVersion string) error {
	fmt.Printf("大版本升级到 %s 暂未实现。\n", targetVersion)
	return fmt.Errorf("大版本升级暂不支持")
}
