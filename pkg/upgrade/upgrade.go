package upgrade

import "fmt"

func Upgrade(instance, namespace, targetVersion string, force, skipChecks bool) error {
	fmt.Printf("Upgrading %s in %s to %s\n", instance, namespace, targetVersion)
	if force {
		fmt.Println("Force upgrade enabled")
	}
	if skipChecks {
		fmt.Println("Skipping pre-upgrade checks")
	}
	return nil
}
