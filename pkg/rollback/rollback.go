package rollback

import "fmt"

func Rollback(instance, namespace string, force bool) error {
	fmt.Printf("Rolling back %s in %s\n", instance, namespace)
	if force {
		fmt.Println("Force rollback enabled")
	}
	return nil
}
