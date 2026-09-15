package rel

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/amankr1098/helm-health/internal/output"
	res "github.com/amankr1098/helm-health/internal/resources"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/release"
	"helm.sh/helm/v4/pkg/storage/driver"
)

// FetchHelmRelease computes the health of a Helm release. It returns a populated
// OutputResult (never nil) alongside an error. Infrastructure failures (Helm init,
// cluster connectivity) are returned as errors; a missing or uninstalled release is
// reported via a StatusNotFound result with a nil error, so callers such as an HTTP
// server can respond gracefully instead of terminating the process.
func FetchHelmRelease(releaseName string, namespace string) (*output.OutputResult, error) {
	startTime := time.Now()
	result := output.NewOutputResult(releaseName, namespace)

	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER")); err != nil {
		return nil, fmt.Errorf("initializing Helm: %w", err)
	}

	releaseGet := action.NewGet(actionConfig)
	rel, err := releaseGet.Run(releaseName)
	if err != nil {
		if errors.Is(err, driver.ErrReleaseNotFound) {
			result.Status = output.StatusNotFound
			result.Message = "release not found — it may have been uninstalled"
			result.Finalize(startTime)
			return result, nil
		}
		return nil, fmt.Errorf("fetching release %q: %w", releaseName, err)
	}

	if rel == nil {
		result.Status = output.StatusNotFound
		result.Message = "release not found — it may have been uninstalled"
		result.Finalize(startTime)
		return result, nil
	}

	releaseResult, err := release.NewAccessor(rel)
	if err != nil {
		return nil, fmt.Errorf("reading release %q: %w", releaseName, err)
	}

	if status := releaseResult.Status(); status != "deployed" {
		if status == "uninstalled" {
			result.Status = output.StatusNotFound
			result.Message = "release has been uninstalled"
		} else {
			result.Status = output.StatusUnhealthy
			result.Message = fmt.Sprintf("release status is %q (expected \"deployed\")", status)
		}
		result.Finalize(startTime)
		return result, nil
	}

	resources, err := processManifest(releaseResult.Manifest(), namespace)
	if err != nil {
		return nil, fmt.Errorf("checking resources for release %q: %w", releaseName, err)
	}
	for _, r := range resources {
		result.AddResource(r)
	}

	result.Finalize(startTime)
	return result, nil
}

// ReleaseInfo is a lightweight summary of a Helm release for listing.
type ReleaseInfo struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Status    string `json:"status"`
	Revision  int    `json:"revision"`
	Updated   string `json:"updated,omitempty"`
}

// ListReleases returns Helm releases across all statuses. If namespace is empty,
// releases from every namespace are returned.
func ListReleases(namespace string) ([]ReleaseInfo, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER")); err != nil {
		return nil, fmt.Errorf("initializing Helm: %w", err)
	}

	list := action.NewList(actionConfig)
	list.All = true
	list.AllNamespaces = namespace == ""
	list.SetStateMask()

	releases, err := list.Run()
	if err != nil {
		return nil, fmt.Errorf("listing releases: %w", err)
	}

	infos := make([]ReleaseInfo, 0, len(releases))
	for _, rel := range releases {
		acc, err := release.NewAccessor(rel)
		if err != nil {
			continue
		}
		info := ReleaseInfo{
			Name:      acc.Name(),
			Namespace: acc.Namespace(),
			Status:    acc.Status(),
			Revision:  acc.Version(),
		}
		if t := acc.DeployedAt(); !t.IsZero() {
			info.Updated = t.UTC().Format(time.RFC3339)
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func processManifest(manifest string, namespace string) ([]output.Resource, error) {
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

	clientset, err := res.GetClientset("")
	if err != nil {
		return nil, err
	}

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
			case "Ingress":
				results = append(results, res.FetchIngress(clientset, namespace, name))
			case "NetworkPolicy":
				results = append(results, res.FetchNetworkPolicy(clientset, namespace, name))
			default:
				// Resources without specific health checks (ConfigMap, Secret, etc.)
				nr := output.NewResource(kind, name, namespace)
				nr.SetStatus(output.StatusHealthy)
				results = append(results, *nr)
			}
		}
	}

	return results, nil
}
