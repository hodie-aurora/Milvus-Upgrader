package k8s

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSet 包含 Kubernetes 客户端和动态客户端
type ClientSet struct {
	KubernetesClient *kubernetes.Clientset
	DynamicClient    dynamic.Interface
}

// GetClient 初始化并返回 Kubernetes 客户端
func GetClient(kubeconfig string) (*ClientSet, error) {
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}
	}

	log.Printf("Using kubeconfig: %s", kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Printf("Failed to build Kubernetes config: %v", err)
		return nil, fmt.Errorf("构建 Kubernetes 配置失败: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create Kubernetes client: %v", err)
		return nil, fmt.Errorf("创建 Kubernetes 客户端失败: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create dynamic client: %v", err)
		return nil, fmt.Errorf("创建动态客户端失败: %v", err)
	}

	log.Printf("Dynamic client: %+v", dynClient)

	return &ClientSet{KubernetesClient: client, DynamicClient: dynClient}, nil
}
