package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/data"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchServices(clientset *kubernetes.Clientset, namespace string, serviceName string) data.ServiceStatus {

	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Service %q in namespace %q: %+v\n", serviceName, namespace, err)
		os.Exit(1)
	}
	endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Endpoints for Service %q in namespace %q: %+v\n", serviceName, namespace, err)
		os.Exit(1)
	}
	return parseService(service, endpoint)
}

func parseService(service *corev1.Service, endpoint *corev1.Endpoints) data.ServiceStatus {
	resourceStatus := data.ServiceStatus{
		Name: service.Name,
		Kind: "Service",
	}
	resourceStatus.Type = string(service.Spec.Type)
	if service.Spec.ClusterIP != "" {
		resourceStatus.Healthy = true
	}
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			resourceStatus.Healthy = true
		} else {
			resourceStatus.Healthy = false
		}
	}

	hasReadyAddresses := false
	allNotReady := true
	if len(endpoint.Subsets) == 0 {
		resourceStatus.Healthy = false
	}

	for _, subset := range endpoint.Subsets {
		if len(subset.Addresses) > 0 {
			hasReadyAddresses = true
			break
		}
	}

	if !hasReadyAddresses {
		for _, subset := range endpoint.Subsets {
			if len(subset.NotReadyAddresses) > 0 {
				allNotReady = false
				break
			}
		}
	}

	if !hasReadyAddresses && allNotReady {
		resourceStatus.Healthy = false
	}
	if hasReadyAddresses {
		resourceStatus.Healthy = true
	}

	return resourceStatus
}
