package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchServices(clientset *kubernetes.Clientset, namespace string, serviceName string) output.Resource {
	service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Service %q in namespace %q: %v\n", serviceName, namespace, err)
		r := output.NewResource("Service", serviceName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	endpoint, err := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), serviceName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Endpoints for Service %q in namespace %q: %v\n", serviceName, namespace, err)
		r := output.NewResource("Service", serviceName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch endpoints: %v", err))
		return *r
	}

	return parseService(service, endpoint)
}

func parseService(service *corev1.Service, endpoint *corev1.Endpoints) output.Resource {
	r := output.NewResource("Service", service.Name, service.Namespace)

	r.SetHealth("type", string(service.Spec.Type))
	r.SetHealth("clusterIP", service.Spec.ClusterIP)

	// Count ready and total endpoints
	readyCount := 0
	notReadyCount := 0
	for _, subset := range endpoint.Subsets {
		readyCount += len(subset.Addresses)
		notReadyCount += len(subset.NotReadyAddresses)
	}
	totalCount := readyCount + notReadyCount

	r.SetHealth("endpoints", map[string]any{
		"ready":    readyCount,
		"notReady": notReadyCount,
		"total":    totalCount,
	})

	healthy := true

	// Check LoadBalancer ingress
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			r.SetHealth("loadBalancer", "ready")
		} else {
			r.SetHealth("loadBalancer", "pending")
			healthy = false
			r.AddIssue("- LoadBalancer ingress not yet assigned")
		}
	}

	// Check endpoints
	if readyCount == 0 && totalCount == 0 {
		// ExternalName or headless services may have no endpoints
		if service.Spec.Type != corev1.ServiceTypeExternalName && service.Spec.ClusterIP != "None" {
			healthy = false
			r.AddIssue("- No endpoints available")
		}
	} else if readyCount < totalCount {
		healthy = false
		r.AddIssue(fmt.Sprintf("- Warning: Only %d of %d expected endpoints available", readyCount, totalCount))
	}

	if healthy {
		r.SetStatus(output.StatusHealthy)
		msg := ""
		if readyCount > 0 {
			msg = fmt.Sprintf("%d endpoints", readyCount)
		}
		if service.Spec.Type == corev1.ServiceTypeLoadBalancer && len(service.Status.LoadBalancer.Ingress) > 0 {
			if msg != "" {
				msg += ", LoadBalancer ready"
			} else {
				msg = "LoadBalancer ready"
			}
		}
		if msg != "" {
			r.SetMessage(msg)
		}
	} else {
		r.SetStatus(output.StatusUnhealthy)
		if totalCount > 0 {
			r.SetMessage(fmt.Sprintf("%d/%d endpoints ready", readyCount, totalCount))
		}
	}

	return *r
}
