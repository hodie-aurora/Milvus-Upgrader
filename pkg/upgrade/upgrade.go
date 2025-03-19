package upgrade

import (
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"github.com/hodie-aurora/milvus-upgrader/pkg/utils"
)

// Upgrade 执行 Milvus 升级
func Upgrade(instance, namespace, sourceVersion, targetVersion string, force, skipChecks bool, kubeconfig string) error {
	client, err := k8s.GetClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("获取 Kubernetes 客户端失败: %v", err)
	}
	sourceVer, err := utils.ParseVersion(sourceVersion)
	if err != nil {
		return fmt.Errorf("源版本解析失败: %v", err)
	}
	targetVer, err := utils.ParseVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("目标版本解析失败: %v", err)
	}
	fmt.Printf("升级 %s (%s) 从 %s 到 %s\n", instance, namespace, sourceVersion, targetVersion)
	if force {
		fmt.Println("强制升级已启用")
	}
	if skipChecks {
		fmt.Println("跳过升级前检查")
	}
	if utils.IsMinorUpgrade(sourceVer, targetVer) {
		return MinorUpgrade(client, instance, namespace, targetVersion)
	}
	return MajorUpgrade(client, instance, namespace, targetVersion)
}
