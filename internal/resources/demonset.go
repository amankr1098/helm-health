package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/data"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchDaemonSet(clientset *kubernetes.Clientset, namespace string, daemonSetName string) data.DaemonSetStatus {

	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	return parseDaemonSet(daemonSet)
}

func parseDaemonSet(daemonSet *appsv1.DaemonSet) data.DaemonSetStatus {
	resourceStatus := data.DaemonSetStatus{
		Name: daemonSet.Name,
		Kind: "DaemonSet",
	}
	resourceStatus.DesiredNumber = daemonSet.Status.DesiredNumberScheduled
	resourceStatus.AvailableNumber = daemonSet.Status.NumberAvailable
	resourceStatus.ReadyNumber = daemonSet.Status.NumberReady
	if resourceStatus.ReadyNumber == resourceStatus.DesiredNumber {
		resourceStatus.Healthy = true
	}
	return resourceStatus
}
