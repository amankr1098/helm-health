package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/output"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchJob(clientset *kubernetes.Clientset, namespace string, jobName string) output.Resource {
	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, v1.GetOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to fetch Job %q in namespace %q: %v\n", jobName, namespace, err)
		r := output.NewResource("Job", jobName, namespace)
		r.SetStatus(output.StatusUnknown)
		r.AddIssue(fmt.Sprintf("- failed to fetch: %v", err))
		return *r
	}

	r := output.NewResource("Job", job.Name, job.Namespace)

	completions := int32(1)
	if job.Spec.Completions != nil {
		completions = *job.Spec.Completions
	}

	r.SetHealth("completions", map[string]any{
		"desired":   completions,
		"succeeded": job.Status.Succeeded,
		"failed":    job.Status.Failed,
		"active":    job.Status.Active,
	})

	if job.Status.Succeeded >= completions {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage(fmt.Sprintf("%d/%d completed", job.Status.Succeeded, completions))
		r.SetHealth("ready", true)
	} else if job.Status.Failed > 0 {
		r.SetStatus(output.StatusUnhealthy)
		r.SetHealth("ready", false)
		r.AddIssue(fmt.Sprintf("- %d failed, %d/%d succeeded", job.Status.Failed, job.Status.Succeeded, completions))
	} else if job.Status.Active > 0 {
		r.SetStatus(output.StatusHealthy)
		r.SetMessage(fmt.Sprintf("%d active, %d/%d completed", job.Status.Active, job.Status.Succeeded, completions))
		r.SetHealth("ready", false)
	} else {
		r.SetStatus(output.StatusUnknown)
		r.SetHealth("ready", false)
	}

	return *r
}
