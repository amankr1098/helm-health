package rel

import (
	"bytes"
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/resources"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/release"
)

func FetchHelmRelease(releaseName string, namespace string) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER")); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	releaseGet := action.NewGet(actionConfig)
	result, err := releaseGet.Run(releaseName)
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	if result != nil {
		releaseResult, err := release.NewAccessor(result)
		if err != nil {
			fmt.Printf("%+v", err)
			os.Exit(1)
		}
		// fmt.Printf("Release: %s, Namespace: %s, Status: %s\n", releaseResult.Name(), releaseResult.Namespace(), releaseResult.Status())

		if releaseResult.Status() == "deployed" {
			fmt.Printf("Release %s is %s.\n", releaseName, releaseResult.Status())
			ProcessManifest(releaseResult.Manifest(), namespace)
		} else {
			fmt.Printf("Release %s is not healthy. Status: %s\n", releaseName, releaseResult.Status())
			os.Exit(1)
		}

	}
}

func ProcessManifest(manifest string, namespace string) {
	// Process the manifest as needed
	// fmt.Printf("Processing manifest:\n%s\n", manifest)
	resourceMap := make(map[string][]string)
	type metaData struct {
		Name string
	}

	type resource struct {
		Kind     string
		Metadata metaData
	}

	manifestByte := bytes.NewReader([]byte(manifest))
	yamlDecoder := yaml.NewDecoder(manifestByte)

	var res resource
	for yamlDecoder.Decode(&res) == nil {
		resourceMap[res.Kind] = append(resourceMap[res.Kind], res.Metadata.Name)
	}

	clientset := resources.GetClientset("")

	for kind, names := range resourceMap {
		// fmt.Printf("Kind: %s, Names: %v\n", kind, names)

		switch kind {
		case "Deployment":
			fmt.Print("\nFetching health for Deployments: \n")
			for _, name := range names {
				result := resources.FetchDeployment(clientset, namespace, name)
				fmt.Printf("Deployment/%s health: %+v\n", name, result)
			}
		case "Service":
			fmt.Print("\nFetching health for Services: \n")
			for _, name := range names {
				result := resources.FetchServices(clientset, namespace, name)
				fmt.Printf("Service/%s health: %+v\n", name, result)
			}
		case "Pod":
			fmt.Print("\nFetching health for Pods: \n")
			for _, name := range names {
				result := resources.FetchPod(clientset, namespace, name)
				fmt.Printf("Pod/%s health: %+v\n", name, result)
			}
		case "PersistentVolumeClaim":
			for _, name := range names {
				fmt.Printf("\nFetching health for PersistentVolumeClaim: %s\n", name)
				result := resources.FetchPVC(clientset, namespace, name)
				fmt.Printf("PersistentVolumeClaim/%s health: %+v\n", name, result)
			}
		case "StatefulSet":
			fmt.Print("\nFetching health for StatefulSets: \n")
			for _, name := range names {
				result := resources.FetchStatefulSet(clientset, namespace, name)
				fmt.Printf("StatefulSet/%s health: %+v\n", name, result)
			}
		case "DaemonSet":
			fmt.Print("\nFetching health for DaemonSets: \n")
			for _, name := range names {
				result := resources.FetchDaemonSet(clientset, namespace, name)
				fmt.Printf("DaemonSet/%s health: %+v\n", name, result)
			}
		case "Job":
			fmt.Print("\nFetching health for Jobs: \n")
			for _, name := range names {
				result := resources.FetchJob(clientset, namespace, name)
				fmt.Printf("Job/%s health: %+v\n", name, result)
			}
		default:
			// fmt.Printf("Resource kind %s is not supported for health checks.\n", kind)
		}
	}

}
