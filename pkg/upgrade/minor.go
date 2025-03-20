package upgrade

import (
	"context"
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	milvusv1beta1 "github.com/milvus-io/milvus-operator/apis/milvus.io/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// MinorUpgrade performs a minor version upgrade
func MinorUpgrade(client *k8s.ClientSet, instance, namespace, targetVersion string) error {
	cr, err := getMilvusCR(client, namespace, instance)
	if err != nil {
		return fmt.Errorf("failed to get Milvus CR: %v", err)
	}
	// Convert to unstructured
	unstructuredCr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cr)
	if err != nil {
		return fmt.Errorf("failed to convert CR to unstructured: %v", err)
	}
	// Dynamically set spec.components.image
	components, found, err := unstructured.NestedMap(unstructuredCr, "spec", "components")
	if !found || err != nil {
		return fmt.Errorf("failed to get spec.components: %v", err)
	}
	components["image"] = "milvusdb/milvus:" + targetVersion
	err = unstructured.SetNestedMap(unstructuredCr, components, "spec", "components")
	if err != nil {
		return fmt.Errorf("failed to set spec.components.image: %v", err)
	}
	// Update
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}
	_, err = client.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), &unstructured.Unstructured{Object: unstructuredCr}, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update Milvus CR: %v", err)
	}
	fmt.Println("Minor version upgrade started, please check cluster status.")
	return nil
}

// getMilvusCR retrieves the Milvus CR
func getMilvusCR(client *k8s.ClientSet, namespace, name string) (*milvusv1beta1.Milvus, error) {
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}

	obj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("error: Milvus instance %s not found in namespace %s", name, namespace)
		}
		return nil, fmt.Errorf("failed to get Milvus CR: %v", err)
	}

	scheme := runtime.NewScheme()
	err = milvusv1beta1.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("failed to add Milvus scheme: %v", err)
	}

	cr := &milvusv1beta1.Milvus{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), cr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Milvus CR: %v", err)
	}

	// Add temporary annotation
	if cr.ObjectMeta.Annotations == nil {
		cr.ObjectMeta.Annotations = make(map[string]string)
	}
	cr.ObjectMeta.Annotations["milvus-upgrader/reconcile-trigger"] = "true"

	// Convert Milvus CR to Unstructured
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cr)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Milvus CR to unstructured: %v", err)
	}

	// Update Milvus CR
	updatedObj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), &unstructured.Unstructured{Object: unstructuredObj}, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to update Milvus CR: %v", err)
	}

	updatedCR := &milvusv1beta1.Milvus{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(updatedObj.UnstructuredContent(), updatedCR)
	if err != nil {
		return nil, fmt.Errorf("failed to convert updated unstructured to Milvus CR: %v", err)
	}

	return updatedCR, nil
}
