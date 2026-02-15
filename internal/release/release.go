package rel

import (
	"fmt"
	"os"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/cli"
)

func GetReleaseHealth(releaseName string, namespace string) {
	// fmt.Printf("Getting health for release: %s in namespace: %s\n", releaseName, namespace)

}

func FetchHelmRelease(releaseName string, namespace string) {
	settings := cli.New()
	actionConfig := new(action.Configuration)

	if err := actionConfig.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER")); err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewList(actionConfig)
	results, err := client.Run()
	if err != nil {
		fmt.Printf("%+v", err)
		os.Exit(1)
	}
	fmt.Print(results)

	// for _, rel := range results {
	// 	fmt.Printf("%+v", rel)
	// }
}
