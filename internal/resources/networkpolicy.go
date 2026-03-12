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

func FetchNetworkPolicy(clientset *kubernetes.Clientset, namespace string, policyName string) output.Resource {
	policy, err := clientset.NetworkingV1().NetworkPolicies(namespace).Get(context.TODO(), policyName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch NetworkPolicy %q in namespace %q: %v\n", policyName, namespace, err)
		r := output.NewResource("NetworkPolicy", policyName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	return parseNetworkPolicy(clientset, policy)
}

func parseNetworkPolicy(clientset *kubernetes.Clientset, policy *networkingv1.NetworkPolicy) output.Resource {
	r := output.NewResource("NetworkPolicy", policy.Name, policy.Namespace)

	healthy := true

	// Parse pod selector
	selector := v1.FormatLabelSelector(&policy.Spec.PodSelector)
	r.SetHealth("podSelector", selector)

	// Count matching pods
	pods, err := clientset.CoreV1().Pods(policy.Namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector,
	})
	matchingPods := 0
	if err != nil {
		r.SetHealth("matchingPods", "unknown")
		r.AddIssue(fmt.Sprintf("- failed to list matching pods: %v", err))
	} else {
		matchingPods = len(pods.Items)
		r.SetHealth("matchingPods", matchingPods)
	}

	// Parse policy types
	policyTypes := []string{}
	for _, pt := range policy.Spec.PolicyTypes {
		policyTypes = append(policyTypes, string(pt))
	}
	r.SetHealth("policyTypes", policyTypes)

	// Parse ingress rules
	ingressRules := []map[string]any{}
	for _, rule := range policy.Spec.Ingress {
		ruleDetail := map[string]any{}
		ports := []string{}
		for _, port := range rule.Ports {
			proto := "TCP"
			if port.Protocol != nil {
				proto = string(*port.Protocol)
			}
			if port.Port != nil {
				ports = append(ports, fmt.Sprintf("%s/%s", proto, port.Port.String()))
			} else {
				ports = append(ports, fmt.Sprintf("%s/*", proto))
			}
		}
		ruleDetail["ports"] = ports

		from := []string{}
		for _, peer := range rule.From {
			if peer.PodSelector != nil {
				from = append(from, fmt.Sprintf("podSelector: %s", v1.FormatLabelSelector(peer.PodSelector)))
			}
			if peer.NamespaceSelector != nil {
				from = append(from, fmt.Sprintf("namespaceSelector: %s", v1.FormatLabelSelector(peer.NamespaceSelector)))
			}
			if peer.IPBlock != nil {
				entry := fmt.Sprintf("ipBlock: %s", peer.IPBlock.CIDR)
				if len(peer.IPBlock.Except) > 0 {
					entry += fmt.Sprintf(" except %v", peer.IPBlock.Except)
				}
				from = append(from, entry)
			}
		}
		ruleDetail["from"] = from
		ingressRules = append(ingressRules, ruleDetail)
	}
	r.SetHealth("ingressRules", ingressRules)

	// Parse egress rules
	egressRules := []map[string]any{}
	for _, rule := range policy.Spec.Egress {
		ruleDetail := map[string]any{}
		ports := []string{}
		for _, port := range rule.Ports {
			proto := "TCP"
			if port.Protocol != nil {
				proto = string(*port.Protocol)
			}
			if port.Port != nil {
				ports = append(ports, fmt.Sprintf("%s/%s", proto, port.Port.String()))
			} else {
				ports = append(ports, fmt.Sprintf("%s/*", proto))
			}
		}
		ruleDetail["ports"] = ports

		to := []string{}
		for _, peer := range rule.To {
			if peer.PodSelector != nil {
				to = append(to, fmt.Sprintf("podSelector: %s", v1.FormatLabelSelector(peer.PodSelector)))
			}
			if peer.NamespaceSelector != nil {
				to = append(to, fmt.Sprintf("namespaceSelector: %s", v1.FormatLabelSelector(peer.NamespaceSelector)))
			}
			if peer.IPBlock != nil {
				entry := fmt.Sprintf("ipBlock: %s", peer.IPBlock.CIDR)
				if len(peer.IPBlock.Except) > 0 {
					entry += fmt.Sprintf(" except %v", peer.IPBlock.Except)
				}
				to = append(to, entry)
			}
		}
		ruleDetail["to"] = to
		egressRules = append(egressRules, ruleDetail)
	}
	r.SetHealth("egressRules", egressRules)

	// Health assessment
	// A NetworkPolicy is considered healthy if it exists and targets at least one pod.
	// If no pods match, it's a warning (policy may be misconfigured or pods not yet deployed).
	if matchingPods == 0 && err == nil {
		healthy = false
		r.AddIssue("- No pods match the policy's pod selector")
	}

	if healthy {
		r.SetStatus(output.StatusHealthy)
		msg := fmt.Sprintf("%d matching pods", matchingPods)
		if len(policyTypes) > 0 {
			msg += fmt.Sprintf(", types: %v", policyTypes)
		}
		r.SetMessage(msg)
	} else {
		r.SetStatus(output.StatusUnhealthy)
		r.SetMessage("no matching pods")
	}

	return *r
}
