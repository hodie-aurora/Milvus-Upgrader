package upgrade

import (
	"context"
	"fmt"

	"github.com/hodie-aurora/milvus-upgrader/pkg/k8s"
	"github.com/hodie-aurora/milvus-upgrader/pkg/utils"
	milvusv1beta1 "github.com/milvus-io/milvus-operator/apis/milvus.io/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// Upgrade performs a Milvus upgrade
func Upgrade(instance, namespace, sourceVersion, targetVersion string, force, skipChecks bool, kubeconfig string) error {
	// Get Kubernetes client
	client, err := k8s.GetClient(kubeconfig)
	if err != nil {
		return fmt.Errorf("failed to get Kubernetes client: %v", err)
	}

	// If sourceVersion is not provided, detect the current version
	if sourceVersion == "" {
		sourceVersion, err = utils.GetCurrentVersion(client, namespace, instance)
		if err != nil {
			return fmt.Errorf("failed to get current version: %v", err)
		}
		fmt.Printf("Detected current version: %s\n", sourceVersion)
	}

	// Parse source version
	sourceVer, err := utils.ParseVersion(sourceVersion)
	if err != nil {
		return fmt.Errorf("failed to parse source version: %v", err)
	}

	// Parse target version
	targetVer, err := utils.ParseVersion(targetVersion)
	if err != nil {
		return fmt.Errorf("failed to parse target version: %v", err)
	}

	// Print upgrade information
	fmt.Printf("Upgrading %s (%s) from %s to %s\n", instance, namespace, sourceVersion, targetVersion)
	if force {
		fmt.Println("Force upgrade enabled")
	}
	if skipChecks {
		fmt.Println("Pre-upgrade checks skipped")
	}

	// Perform minor or major upgrade based on version comparison
	if utils.IsMinorUpgrade(sourceVer, targetVer) {
		return MinorUpgrade(client, instance, namespace, targetVersion)
	}
	return MajorUpgrade(client, instance, namespace, targetVersion)
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
