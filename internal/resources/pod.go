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

func FetchPod(clientset *kubernetes.Clientset, namespace string, podName string) data.PodStatus {

	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	return parsePod(pod)
}

func parsePod(pod *corev1.Pod) data.PodStatus {
	resourceStatus := data.PodStatus{
		Name: pod.Name,
		Kind: "Pod",
	}
	resourceStatus.Phase = string(pod.Status.Phase)
	allReady := true
	for _, containerStatus := range pod.Status.ContainerStatuses {
		if !containerStatus.Ready {
			allReady = false
			break
		}
	}

	resourceStatus.Ready = allReady
	resourceStatus.Containers = len(pod.Spec.Containers)
	if allReady && pod.Status.Phase == corev1.PodRunning {
		resourceStatus.Healthy = true
	} else {
		resourceStatus.Healthy = false
	}

	return resourceStatus
}
