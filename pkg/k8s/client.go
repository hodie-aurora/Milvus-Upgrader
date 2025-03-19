package k8s

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClientSet contains Kubernetes client and dynamic client
type ClientSet struct {
	KubernetesClient *kubernetes.Clientset
	DynamicClient    dynamic.Interface
}

// GetClient initializes and returns a Kubernetes client
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
		return nil, fmt.Errorf("Failed to build Kubernetes config: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create Kubernetes client: %v", err)
		return nil, fmt.Errorf("Failed to create Kubernetes client: %v", err)
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Printf("Failed to create dynamic client: %v", err)
		return nil, fmt.Errorf("Failed to create dynamic client: %v", err)
	}

	log.Printf("Dynamic client: %+v", dynClient)

	return &ClientSet{KubernetesClient: client, DynamicClient: dynClient}, nil
}
