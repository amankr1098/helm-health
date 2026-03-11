package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchPVC(clientset *kubernetes.Clientset, namespace string, pvcName string) output.Resource {
	pvc, err := clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch PVC %q in namespace %q: %v\n", pvcName, namespace, err)
		r := output.NewResource("PersistentVolumeClaim", pvcName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	r := output.NewResource("PersistentVolumeClaim", pvc.Name, pvc.Namespace)

	phase := string(pvc.Status.Phase)
	r.SetHealth("phase", phase)

	if phase == "Bound" {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage("Bound")
	} else {
		r.SetStatus(output.StatusUnhealthy)
		r.AddIssue(fmt.Sprintf("- Phase: %s", phase))
	}

	return *r
}
