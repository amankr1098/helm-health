package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchStatefulSet(clientset *kubernetes.Clientset, namespace string, statefulSetName string) output.Resource {
	statefulSet, err := clientset.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch StatefulSet %q in namespace %q: %v\n", statefulSetName, namespace, err)
		r := output.NewResource("StatefulSet", statefulSetName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	r := output.NewResource("StatefulSet", statefulSet.Name, statefulSet.Namespace)

	desired := statefulSet.Status.Replicas
	ready := statefulSet.Status.ReadyReplicas
	current := statefulSet.Status.CurrentReplicas
	updated := statefulSet.Status.UpdatedReplicas

	r.SetHealth("replicas", map[string]any{
		"desired": desired,
		"ready":   ready,
		"current": current,
		"updated": updated,
	})

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
		r.AddIssue(fmt.Sprintf("- %d/%d replicas ready", ready, desired))
	}

	return *r
}
