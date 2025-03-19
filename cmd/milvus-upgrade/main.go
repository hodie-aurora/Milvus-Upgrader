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
		Short: "Upgrade Milvus to the target version",
		Run: func(cmd *cobra.Command, args []string) {
			err := upgrade.Upgrade(instance, namespace, sourceVersion, targetVersion, force, skipChecks, kubeconfig)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVarP(&instance, "instance", "i", "", "Milvus instance name (required)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().StringVarP(&sourceVersion, "source-version", "s", "", "Current Milvus version (required)")
	cmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "Target Milvus version (required)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force upgrade without confirmation")
	cmd.Flags().BoolVarP(&skipChecks, "skip-checks", "k", false, "Skip pre-upgrade checks")
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig file")
	cmd.MarkFlagRequired("instance")
	cmd.MarkFlagRequired("source-version")
	cmd.MarkFlagRequired("target-version")
	return cmd
}
