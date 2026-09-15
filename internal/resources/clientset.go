package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// RequestTimeout bounds each Kubernetes API interaction so a health check can never
// hang indefinitely against an unreachable cluster.
const RequestTimeout = 15 * time.Second

// RequestContext returns a context carrying the default per-request timeout.
func RequestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), RequestTimeout)
}

func GetClientset(kubeconfigPath string) (*kubernetes.Clientset, error) {
	if kubeconfigPath == "" {
		kubeconfigPath = GetKubeconfigPath("")
	}
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("building kube config from %q: %w", kubeconfigPath, err)
	}

	// Create a Kubernetes client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating kubernetes client: %w", err)
	}
	return clientset, nil
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
