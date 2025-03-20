package upgrade

import (
	"context"
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"github.com/hodie-aurora/milvus-upgrader/pkg/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// MajorUpgrade performs a major version upgrade
func MajorUpgrade(client *k8s.ClientSet, instance, namespace, targetVersion string) error {
	// Check dependency versions
	err := utils.CheckDependencies(client, namespace, instance, targetVersion)
	if err != nil {
		return fmt.Errorf("major version upgrade aborted: %v", err)
	}

	// Get the current Milvus CR
	cr, err := getMilvusCR(client, namespace, instance)
	if err != nil {
		return fmt.Errorf("failed to get Milvus CR: %v", err)
	}

	// Convert to unstructured
	unstructuredCr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cr)
	if err != nil {
		return fmt.Errorf("failed to convert CR to unstructured: %v", err)
	}

	// Dynamically set spec.components.image and enable rolling update
	components, found, err := unstructured.NestedMap(unstructuredCr, "spec", "components")
	if !found || err != nil {
		return fmt.Errorf("failed to get spec.components: %v", err)
	}
	components["image"] = "milvusdb/milvus:" + targetVersion
	components["enableRollingUpdate"] = true
	components["imageUpdateMode"] = "rollingUpgrade"
	err = unstructured.SetNestedMap(unstructuredCr, components, "spec", "components")
	if err != nil {
		return fmt.Errorf("failed to set spec.components: %v", err)
	}

	// Update the Milvus CR
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}
	_, err = client.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), &unstructured.Unstructured{Object: unstructuredCr}, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Milvus CR: %v", err)
	}

	fmt.Printf("Major version upgrade to %s started. Rolling upgrade in progress, please check cluster status.\n", targetVersion)
	return nil
}
