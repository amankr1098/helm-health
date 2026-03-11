package resources

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/amankr1098/helm-health/internal/output"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) output.Resource {
	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Deployment %q in namespace %q: %v\n", deploymentName, namespace, err)
		r := output.NewResource("Deployment", deploymentName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}
	return parseDeployment(clientset, deployment)
}

func parseDeployment(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) output.Resource {
	r := output.NewResource("Deployment", deployment.Name, deployment.Namespace)

	desired := int32(0)
	if deployment.Spec.Replicas != nil {
		desired = *deployment.Spec.Replicas
	}
	ready := deployment.Status.ReadyReplicas
	available := deployment.Status.AvailableReplicas
	unavailable := deployment.Status.UnavailableReplicas

	replicas := map[string]any{
		"desired":     desired,
		"ready":       ready,
		"available":   available,
		"unavailable": unavailable,
	}

	// Deployment conditions for JSON output
	conditions := []map[string]any{}
	for _, c := range deployment.Status.Conditions {
		conditions = append(conditions, map[string]any{
			"type":    string(c.Type),
			"status":  string(c.Status),
			"reason":  c.Reason,
			"message": c.Message,
		})
	}

	r.SetHealth("replicas", replicas)
	r.SetHealth("conditions", conditions)

	if ready >= desired && desired > 0 {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage(fmt.Sprintf("%d/%d replicas ready", ready, desired))
		r.SetHealth("ready", true)
	} else if desired == 0 {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage("scaled to 0")
		r.SetHealth("ready", true)
	} else {
		r.SetStatus(output.StatusUnhealthy)
		r.SetHealth("ready", false)
		r.AddIssue(fmt.Sprintf("- %d/%d replicas available", available, desired))

		// Fetch pods for detailed diagnostics
		pods := fetchDeploymentPods(clientset, deployment)
		podDetails := []map[string]any{}
		for i := range pods {
			pod := &pods[i]
			podDetail := BuildPodHealthMap(pod)
			podDetails = append(podDetails, podDetail)
			for _, line := range PodIssueLines(pod) {
				r.AddIssue(line)
			}
		}
		if len(podDetails) > 0 {
			r.SetHealth("pods", podDetails)
		}
	}

	return *r
}

func fetchDeploymentPods(clientset *kubernetes.Clientset, deployment *appsv1.Deployment) []corev1.Pod {
	selector := deployment.Spec.Selector
	if selector == nil {
		return nil
	}
	labelSelector := v1.FormatLabelSelector(selector)
	pods, err := clientset.CoreV1().Pods(deployment.Namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil
	}
	return pods.Items
}

// BuildPodHealthMap creates a JSON-friendly map of pod health details.
func BuildPodHealthMap(pod *corev1.Pod) map[string]any {
	containers := []map[string]any{}
	podHealthy := true
	podReason := ""

	for _, cs := range pod.Status.ContainerStatuses {
		c := map[string]any{
			"name":         cs.Name,
			"ready":        cs.Ready,
			"restartCount": cs.RestartCount,
		}
		if cs.State.Waiting != nil {
			c["state"] = "waiting"
			c["reason"] = cs.State.Waiting.Reason
			c["message"] = cs.State.Waiting.Message
			podHealthy = false
			podReason = cs.State.Waiting.Reason
		} else if cs.State.Terminated != nil {
			c["state"] = "terminated"
			c["reason"] = cs.State.Terminated.Reason
			c["message"] = cs.State.Terminated.Message
			c["exitCode"] = cs.State.Terminated.ExitCode
			podHealthy = false
			podReason = cs.State.Terminated.Reason
		} else if cs.State.Running != nil {
			c["state"] = "running"
		}
		if cs.LastTerminationState.Terminated != nil {
			c["lastRestart"] = cs.LastTerminationState.Terminated.FinishedAt.Format(time.RFC3339)
		}
		if !cs.Ready {
			podHealthy = false
		}
		containers = append(containers, c)
	}

	status := output.StatusHealthy
	if !podHealthy {
		status = output.StatusUnhealthy
	}

	return map[string]any{
		"name":       pod.Name,
		"status":     status,
		"phase":      string(pod.Status.Phase),
		"reason":     podReason,
		"containers": containers,
	}
}

// PodIssueLines builds indented text lines for unhealthy pod containers.
func PodIssueLines(pod *corev1.Pod) []string {
	var lines []string
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Waiting != nil {
			lines = append(lines, fmt.Sprintf("- Pod %s: %s", pod.Name, cs.State.Waiting.Reason))
			lines = append(lines, fmt.Sprintf("  Container: %s", cs.Name))
			lines = append(lines, fmt.Sprintf("  Reason: %s", cs.State.Waiting.Reason))
			if cs.State.Waiting.Message != "" {
				lines = append(lines, fmt.Sprintf("  Message: %s", cs.State.Waiting.Message))
			}
			if cs.LastTerminationState.Terminated != nil {
				elapsed := time.Since(cs.LastTerminationState.Terminated.FinishedAt.Time)
				lines = append(lines, fmt.Sprintf("  Last restart: %s ago", FormatDuration(elapsed)))
			}
		} else if cs.State.Terminated != nil && cs.State.Terminated.ExitCode != 0 {
			lines = append(lines, fmt.Sprintf("- Pod %s: %s", pod.Name, cs.State.Terminated.Reason))
			lines = append(lines, fmt.Sprintf("  Container: %s", cs.Name))
			lines = append(lines, fmt.Sprintf("  Reason: %s", cs.State.Terminated.Reason))
			if cs.State.Terminated.Message != "" {
				lines = append(lines, fmt.Sprintf("  Message: %s", cs.State.Terminated.Message))
			}
		}
	}
	return lines
}

// FormatDuration returns a human-friendly duration string.
func FormatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	return fmt.Sprintf("%dh", int(d.Hours()))
}
