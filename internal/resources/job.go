package resources

import (
	"context"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/data"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func FetchJob(clientset *kubernetes.Clientset, namespace string, jobName string) data.JobStatus {

	job, err := clientset.BatchV1().Jobs(namespace).Get(context.TODO(), jobName, v1.GetOptions{})
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	return parseJob(job)
}

func parseJob(job *batchv1.Job) data.JobStatus {
	resourceStatus := data.JobStatus{
		Name: job.Name,
		Kind: "Job",
	}
	resourceStatus.Succeeded = job.Status.Succeeded
	resourceStatus.Failed = job.Status.Failed
	resourceStatus.Active = job.Status.Active
	resourceStatus.Completions = *job.Spec.Completions
	if resourceStatus.Succeeded >= *job.Spec.Completions {
		resourceStatus.Healthy = true
	}
	return resourceStatus
}
