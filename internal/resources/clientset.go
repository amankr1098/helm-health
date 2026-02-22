package resources

import (
	"fmt"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientset(kubeconfigPath string) *kubernetes.Clientset {
	if kubeconfigPath == "" {
		kubeconfigPath = GetKubeconfigPath("")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	// Create a Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	return clientset
}

func GetKubeconfigPath(path string) string {
	if path != "" {
		return path
	}
	path = os.Getenv("KUBECONFIG")
	if path == "" {
		path = os.ExpandEnv("$HOME/.kube/config")
	}
	return path
}
