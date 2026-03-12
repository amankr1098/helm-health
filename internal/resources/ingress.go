package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchIngress(clientset *kubernetes.Clientset, namespace string, ingressName string) output.Resource {
	ingress, err := clientset.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Ingress %q in namespace %q: %v\n", ingressName, namespace, err)
		r := output.NewResource("Ingress", ingressName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	return parseIngress(clientset, ingress)
}

func parseIngress(clientset *kubernetes.Clientset, ingress *networkingv1.Ingress) output.Resource {
	r := output.NewResource("Ingress", ingress.Name, ingress.Namespace)

	healthy := true

	// Check if LoadBalancer ingress IPs/hostnames are assigned
	lbIngress := ingress.Status.LoadBalancer.Ingress
	lbDetails := []map[string]any{}
	for _, lb := range lbIngress {
		entry := map[string]any{}
		if lb.IP != "" {
			entry["ip"] = lb.IP
		}
		if lb.Hostname != "" {
			entry["hostname"] = lb.Hostname
		}
		lbDetails = append(lbDetails, entry)
	}
	r.SetHealth("loadBalancer", lbDetails)

	if len(lbIngress) == 0 {
		healthy = false
		r.AddIssue("- No LoadBalancer ingress IP/hostname assigned")
	}

	// Collect ingress class if set
	if ingress.Spec.IngressClassName != nil {
		r.SetHealth("ingressClass", *ingress.Spec.IngressClassName)
	}

	// Check TLS configuration
	tlsHosts := []string{}
	for _, tls := range ingress.Spec.TLS {
		tlsHosts = append(tlsHosts, tls.Hosts...)
	}
	if len(tlsHosts) > 0 {
		r.SetHealth("tlsHosts", tlsHosts)
	}

	// Check backend services for each rule
	totalBackends := 0
	healthyBackends := 0
	backendDetails := []map[string]any{}

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			totalBackends++
			backend := map[string]any{
				"host": rule.Host,
				"path": path.Path,
			}

			if path.Backend.Service != nil {
				svcName := path.Backend.Service.Name
				backend["service"] = svcName

				// Check if backend service has ready endpoints
				endpoints, err := clientset.CoreV1().Endpoints(ingress.Namespace).Get(context.TODO(), svcName, v1.GetOptions{})
				if err != nil {
					backend["status"] = "unknown"
					backend["error"] = fmt.Sprintf("failed to fetch endpoints: %v", err)
					healthy = false
					r.AddIssue(fmt.Sprintf("- Backend %s: failed to fetch endpoints", svcName))
				} else {
					readyCount := 0
					for _, subset := range endpoints.Subsets {
						readyCount += len(subset.Addresses)
					}
					backend["readyEndpoints"] = readyCount
					if readyCount > 0 {
						backend["status"] = "healthy"
						healthyBackends++
					} else {
						backend["status"] = "noEndpoints"
						healthy = false
						r.AddIssue(fmt.Sprintf("- Backend %s: no ready endpoints", svcName))
					}
				}
			}

			backendDetails = append(backendDetails, backend)
		}
	}

	// Check default backend
	if ingress.Spec.DefaultBackend != nil && ingress.Spec.DefaultBackend.Service != nil {
		totalBackends++
		svcName := ingress.Spec.DefaultBackend.Service.Name
		backend := map[string]any{
			"service": svcName,
			"default": true,
		}
		endpoints, err := clientset.CoreV1().Endpoints(ingress.Namespace).Get(context.TODO(), svcName, v1.GetOptions{})
		if err != nil {
			backend["status"] = "unknown"
			healthy = false
			r.AddIssue(fmt.Sprintf("- Default backend %s: failed to fetch endpoints", svcName))
		} else {
			readyCount := 0
			for _, subset := range endpoints.Subsets {
				readyCount += len(subset.Addresses)
			}
			backend["readyEndpoints"] = readyCount
			if readyCount > 0 {
				backend["status"] = "healthy"
				healthyBackends++
			} else {
				backend["status"] = "noEndpoints"
				healthy = false
				r.AddIssue(fmt.Sprintf("- Default backend %s: no ready endpoints", svcName))
			}
		}
		backendDetails = append(backendDetails, backend)
	}

	r.SetHealth("backends", backendDetails)
	r.SetHealth("totalBackends", totalBackends)
	r.SetHealth("healthyBackends", healthyBackends)

	if healthy {
		r.SetStatus(output.StatusHealthy)
		msg := fmt.Sprintf("%d/%d backends ready", healthyBackends, totalBackends)
		if len(lbIngress) > 0 {
			msg += ", LB assigned"
		}
		r.SetMessage(msg)
	} else if healthyBackends > 0 && healthyBackends < totalBackends {
		r.SetStatus(output.StatusUnhealthy)
		r.SetMessage(fmt.Sprintf("%d/%d backends ready", healthyBackends, totalBackends))
	} else {
		r.SetStatus(output.StatusUnhealthy)
	}

	return *r
}
