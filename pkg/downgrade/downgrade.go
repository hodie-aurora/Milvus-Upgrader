package downgrade

import "fmt"

func Downgrade(instance, namespace, targetVersion string, force, skipChecks bool) error {
	fmt.Printf("Downgrading %s in %s to %s\n", instance, namespace, targetVersion)
	if force {
		fmt.Println("Force downgrade enabled")
	}
	if skipChecks {
		fmt.Println("Skipping pre-downgrade checks")
	}
	return nil
}
