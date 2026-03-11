package rel

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/amankr1098/helm-health/internal/output"
	res "github.com/amankr1098/helm-health/internal/resources"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/release"
)

func FetchHelmRelease(releaseName string, namespace string, format output.OutputFormat) {
	startTime := time.Now()
	result := output.NewOutputResult(releaseName, namespace)

	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER")); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing Helm: %v\n", err)
		os.Exit(1)
	}

	releaseGet := action.NewGet(actionConfig)
	rel, err := releaseGet.Run(releaseName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching release: %v\n", err)
		os.Exit(1)
	}

	if rel == nil {
		fmt.Fprintf(os.Stderr, "Release %q not found\n", releaseName)
		os.Exit(1)
	}

	releaseResult, err := release.NewAccessor(rel)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading release: %v\n", err)
		os.Exit(1)
	}

	if releaseResult.Status() != "deployed" {
		result.Status = output.StatusUnhealthy
		result.Finalize(startTime)
		result.Print(format)
		os.Exit(1)
	}

	resources := processManifest(releaseResult.Manifest(), namespace)
	for _, r := range resources {
		result.AddResource(r)
	}

	result.Finalize(startTime)
	result.Print(format)

	if result.Status == output.StatusUnhealthy {
		os.Exit(1)
	}
}

func processManifest(manifest string, namespace string) []output.Resource {
	type metaData struct {
		Name string
	}
	type resource struct {
		Kind     string
		Metadata metaData
	}

	resourceMap := make(map[string][]string)
	manifestByte := bytes.NewReader([]byte(manifest))
	yamlDecoder := yaml.NewDecoder(manifestByte)

	var r resource
	for yamlDecoder.Decode(&r) == nil {
		resourceMap[r.Kind] = append(resourceMap[r.Kind], r.Metadata.Name)
	}

	clientset := res.GetClientset("")

	var results []output.Resource

	for kind, names := range resourceMap {
		for _, name := range names {
			switch kind {
			case "Deployment":
				results = append(results, res.FetchDeployment(clientset, namespace, name))
			case "StatefulSet":
				results = append(results, res.FetchStatefulSet(clientset, namespace, name))
			case "DaemonSet":
				results = append(results, res.FetchDaemonSet(clientset, namespace, name))
			case "Service":
				results = append(results, res.FetchServices(clientset, namespace, name))
			case "Pod":
				results = append(results, res.FetchPod(clientset, namespace, name))
			case "PersistentVolumeClaim":
				results = append(results, res.FetchPVC(clientset, namespace, name))
			case "Job":
				results = append(results, res.FetchJob(clientset, namespace, name))
			default:
				// Resources without specific health checks (ConfigMap, Secret, etc.)
				nr := output.NewResource(kind, name, namespace)
				nr.SetStatus(output.StatusHealthy)
				results = append(results, *nr)
			}
		}
	}

	return results
}
