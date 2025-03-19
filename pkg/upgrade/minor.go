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

// MinorUpgrade 执行小版本升级
func MinorUpgrade(client *k8s.ClientSet, instance, namespace, targetVersion string) error {
	cr, err := getMilvusCR(client, namespace, instance)
	if err != nil {
		return fmt.Errorf("获取 Milvus CR 失败: %v", err)
	}
	// 转成 unstructured
	unstructuredCr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cr)
	if err != nil {
		return fmt.Errorf("转换 CR 为 Unstructured 失败: %v", err)
	}
	// 动态设置 spec.components.image
	components, found, err := unstructured.NestedMap(unstructuredCr, "spec", "components")
	if !found || err != nil {
		return fmt.Errorf("获取 spec.components 失败: %v", err)
	}
	components["image"] = "milvusdb/milvus:" + targetVersion
	err = unstructured.SetNestedMap(unstructuredCr, components, "spec", "components")
	if err != nil {
		return fmt.Errorf("设置 spec.components.image 失败: %v", err)
	}
	// 更新
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses", // 已修正为 milvuses
	}
	_, err = client.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), &unstructured.Unstructured{Object: unstructuredCr}, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("更新 Milvus CR 失败: %v", err)
	}
	fmt.Println("小版本升级已启动，请检查集群状态。")
	return nil
}

// getMilvusCR 获取 Milvus CR
func getMilvusCR(client *k8s.ClientSet, namespace, name string) (*milvusv1beta1.Milvus, error) {
	gvr := schema.GroupVersionResource{
		Group:    "milvus.io",
		Version:  "v1beta1",
		Resource: "milvuses",
	}

	obj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("Milvus 实例 %s 在命名空间 %s 中未找到", name, namespace)
		}
		return nil, fmt.Errorf("获取 Milvus CR 失败: %v", err)
	}

	scheme := runtime.NewScheme()
	err = milvusv1beta1.AddToScheme(scheme)
	if err != nil {
		return nil, fmt.Errorf("添加 Milvus Scheme 失败: %v", err)
	}

	cr := &milvusv1beta1.Milvus{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), cr)
	if err != nil {
		return nil, fmt.Errorf("转换 Milvus CR 失败: %v", err)
	}

	// 添加临时 annotation
	if cr.ObjectMeta.Annotations == nil {
		cr.ObjectMeta.Annotations = make(map[string]string)
	}
	cr.ObjectMeta.Annotations["milvus-upgrader/reconcile-trigger"] = "true"

	// 将 Milvus CR 转换为 Unstructured
	unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cr)
	if err != nil {
		return nil, fmt.Errorf("转换 Milvus CR 到 Unstructured 失败: %v", err)
	}

	// 更新 Milvus CR
	updatedObj, err := client.DynamicClient.Resource(gvr).Namespace(namespace).Update(context.TODO(), &unstructured.Unstructured{Object: unstructuredObj}, metav1.UpdateOptions{})
	if err != nil {
		return nil, fmt.Errorf("更新 Milvus CR 失败: %v", err)
	}

	updatedCR := &milvusv1beta1.Milvus{}
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(updatedObj.UnstructuredContent(), updatedCR)
	if err != nil {
		return nil, fmt.Errorf("转换更新后的 Unstructured 到 Milvus CR 失败: %v", err)
	}

	return updatedCR, nil
}
