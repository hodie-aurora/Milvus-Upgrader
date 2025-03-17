package main

import (
	"fmt"
	"os"

	"github.com/hodie-aurora/milvus-upgrader/pkg/downgrade"
	"github.com/hodie-aurora/milvus-upgrader/pkg/rollback"
	"github.com/hodie-aurora/milvus-upgrader/pkg/upgrade"
	"github.com/spf13/cobra"
)

var (
	instance      string
	namespace     string
	targetVersion string
	force         bool
	skipChecks    bool
)

func main() {
	var rootCmd = &cobra.Command{Use: "milvus-upgrade"}

	rootCmd.AddCommand(upgradeCmd())
	rootCmd.AddCommand(downgradeCmd())
	rootCmd.AddCommand(rollbackCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func upgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Milvus to a target version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := upgrade.Upgrade(instance, namespace, targetVersion, force, skipChecks); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&instance, "instance", "i", "", "Milvus instance name (required)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "Target Milvus version (required)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force upgrade without confirmation")
	cmd.Flags().BoolVarP(&skipChecks, "skip-checks", "s", false, "Skip pre-upgrade checks")

	cmd.MarkFlagRequired("instance")
	cmd.MarkFlagRequired("target-version")
	return cmd
}

func downgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "downgrade",
		Short: "Downgrade Milvus to a target version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := downgrade.Downgrade(instance, namespace, targetVersion, force, skipChecks); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&instance, "instance", "i", "", "Milvus instance name (required)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().StringVarP(&targetVersion, "target-version", "t", "", "Target Milvus version (required)")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force downgrade without confirmation")
	cmd.Flags().BoolVarP(&skipChecks, "skip-checks", "s", false, "Skip pre-downgrade checks")

	cmd.MarkFlagRequired("instance")
	cmd.MarkFlagRequired("target-version")
	return cmd
}

func rollbackCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollback",
		Short: "Rollback Milvus to the previous version",
		Run: func(cmd *cobra.Command, args []string) {
			if err := rollback.Rollback(instance, namespace, force); err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVarP(&instance, "instance", "i", "", "Milvus instance name (required)")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "Force rollback without confirmation")

	cmd.MarkFlagRequired("instance")
	return cmd
}
