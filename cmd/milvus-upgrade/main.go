package main

import (
	"fmt"
	"os"

	"github.com/hodie-aurora/milvus-upgrader/pkg/upgrade"
	"github.com/spf13/cobra"
)

var (
	instance      string
	namespace     string
	sourceVersion string
	targetVersion string
	force         bool
	skipChecks    bool
	kubeconfig    string
)

func main() {
	rootCmd := &cobra.Command{Use: "milvus-upgrade"}
	rootCmd.AddCommand(upgradeCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "升级 Milvus 到目标版本",
		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade.Upgrade(instance, namespace, sourceVersion, targetVersion, force, skipChecks, kubeconfig)
			if err != nil {
				fmt.Println("错误:", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVarP(&instance, "instance", "i", "", "Milvus 实例名 (必填)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes 命名空间")
	cmd.Flags().StringVarP(&sourceVersion, "source-version", "s", "", "当前 Milvus 版本 (必填)")
	cmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "目标 Milvus 版本 (必填)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "强制升级，不需确认")
	cmd.Flags().BoolVarP(&skipChecks, "skip-checks", "k", false, "跳过升级前检查")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig 文件路径")
	cmd.MarkFlagRequired("instance")
	cmd.MarkFlagRequired("source-version")
	cmd.MarkFlagRequired("target-version")
	return cmd
}
