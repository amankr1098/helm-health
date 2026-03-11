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

func FetchPod(clientset *kubernetes.Clientset, namespace string, podName string) output.Resource {
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Pod %q in namespace %q: %v\n", podName, namespace, err)
		r := output.NewResource("Pod", podName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}
	return parsePod(pod)
}

func parsePod(pod *corev1.Pod) output.Resource {
	r := output.NewResource("Pod", pod.Name, pod.Namespace)

	podHealth := BuildPodHealthMap(pod)
	r.SetHealth("phase", podHealth["phase"])
	r.SetHealth("containers", podHealth["containers"])

	allReady := true
	for _, cs := range pod.Status.ContainerStatuses {
		if !cs.Ready {
			allReady = false
			break
		}
	}

	if allReady && pod.Status.Phase == corev1.PodRunning {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage(fmt.Sprintf("%d/%d containers ready", len(pod.Status.ContainerStatuses), len(pod.Spec.Containers)))
		r.SetHealth("ready", true)
	} else if pod.Status.Phase == corev1.PodSucceeded {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage("Completed")
		r.SetHealth("ready", true)
	} else {
		r.SetStatus(output.StatusUnhealthy)
		r.SetHealth("ready", false)
		r.AddIssue(fmt.Sprintf("- Phase: %s", pod.Status.Phase))
		for _, line := range PodIssueLines(pod) {
			r.AddIssue(line)
		}
	}

	return *r
}
