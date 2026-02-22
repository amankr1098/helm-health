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

func FetchStatefulSet(clientset *kubernetes.Clientset, namespace string, statefulSetName string) data.StatefulSetStatus {

	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("Failed to fetch StatefulSet %q in namespace %q: %v\n", statefulSetName, namespace, err)
		os.Exit(1)
	}
	return parseStatefulSet(statefulSet)
}

func parseStatefulSet(statefulSet *appsv1.StatefulSet) data.StatefulSetStatus {
	resourceStatus := data.StatefulSetStatus{
		Name: statefulSet.Name,
		Kind: "StatefulSet",
	}

	resourceStatus.CurrentReplicas = statefulSet.Status.CurrentReplicas
	resourceStatus.DesiredReplicas = statefulSet.Status.Replicas
	resourceStatus.ReadyReplicas = statefulSet.Status.ReadyReplicas
	resourceStatus.UpdatedReplicas = statefulSet.Status.UpdatedReplicas

	if resourceStatus.ReadyReplicas == resourceStatus.DesiredReplicas {
		resourceStatus.Healthy = true
	} else {
		resourceStatus.Healthy = false
	}

	return resourceStatus

}
