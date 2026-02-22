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

func FetchDeployment(clientset *kubernetes.Clientset, namespace string, deploymentName string) data.DeploymentStatus {

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	// fmt.Printf("Deployment Name: %s, Namespace: %s, Replicas: %d\n", deployment.Name, deployment.Namespace, *deployment.Spec.Replicas)

	return parseDeployment(deployment)
}

func parseDeployment(deployment *appsv1.Deployment) data.DeploymentStatus {
	resourceStatus := data.DeploymentStatus{
		Name: deployment.Name,
		Kind: "Deployment",
	}

	desired := int32(0)
	if deployment.Spec.Replicas != nil {
		desired = *deployment.Spec.Replicas
	}

	ready := deployment.Status.ReadyReplicas

	resourceStatus.DesiredReplicas = desired
	resourceStatus.ReadyReplicas = ready

	if ready < desired {
		resourceStatus.Healthy = false
	} else if ready == desired || ready > desired {
		resourceStatus.Healthy = true
	}

	return resourceStatus
}
