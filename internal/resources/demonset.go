package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchDaemonSet(clientset *kubernetes.Clientset, namespace string, daemonSetName string) output.Resource {
	daemonSet, err := clientset.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch DaemonSet %q in namespace %q: %v\n", daemonSetName, namespace, err)
		r := output.NewResource("DaemonSet", daemonSetName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	r := output.NewResource("DaemonSet", daemonSet.Name, daemonSet.Namespace)

	desired := daemonSet.Status.DesiredNumberScheduled
	ready := daemonSet.Status.NumberReady
	available := daemonSet.Status.NumberAvailable
	misscheduled := daemonSet.Status.NumberMisscheduled

	r.SetHealth("replicas", map[string]any{
		"desired":      desired,
		"ready":        ready,
		"available":    available,
		"misscheduled": misscheduled,
	})

	if ready >= desired && desired > 0 {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage(fmt.Sprintf("%d/%d nodes ready", ready, desired))
		r.SetHealth("ready", true)
	} else if desired == 0 {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage("no nodes scheduled")
		r.SetHealth("ready", true)
	} else {
		r.SetStatus(output.StatusUnhealthy)
		r.SetHealth("ready", false)
		r.AddIssue(fmt.Sprintf("- %d/%d nodes ready", ready, desired))
		if misscheduled > 0 {
			r.AddIssue(fmt.Sprintf("- %d nodes misscheduled", misscheduled))
		}
	}

	return *r
}
