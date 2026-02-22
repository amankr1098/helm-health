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

func FetchPVC(clientset *kubernetes.Clientset, namespace string, pvcName string) data.PVCStatus {

	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	return parsePVC(pvc)
}

func parsePVC(pvc *corev1.PersistentVolumeClaim) data.PVCStatus {
	resourceStatus := data.PVCStatus{
		Name: pvc.Name,
		Kind: "PersistentVolumeClaim",
	}
	resourceStatus.Phase = string(pvc.Status.Phase)
	if resourceStatus.Phase == "Bound" {
		resourceStatus.Healthy = true
	}
	return resourceStatus
}
