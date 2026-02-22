/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	rel "github.com/amankr1098/helm-health/internal/release"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "helm-health",
	Short: "Helm-health is a tool to check the health of your Helm releases",
	Long: `Helm-health is a tool to check the health of your Helm releases. 
It checks the status of your releases and provides a report on their health.
It can be used to check the health of your releases in a Kubernetes cluster 
and provide insights into their status.`,
	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		releaseFlag := cmd.Flag("release_name")
		releaseName := releaseFlag.Value.String()
		if releaseName == "" {
			cmd.PrintErrln("Error: --release_name flag is required")
			cmd.Help()
			os.Exit(1)
		}
		// Use HELM_NAMESPACE env var (set by helm when -n is used) if available,
		// otherwise fall back to the flag value
		namespace := os.Getenv("HELM_NAMESPACE")
		if namespace == "" {
			namespaceFlag := cmd.Flag("namespace")
			namespace = namespaceFlag.Value.String()
		}
		fmt.Printf("Release Name: %s, Namespace: %s\n", releaseName, namespace)
		rel.FetchHelmRelease(releaseName, namespace)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.helm-health.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.Flags().StringP("release_name", "r", "", "Name of the Helm release to check")
	rootCmd.Flags().StringP("namespace", "n", "default", "Namespace of the Helm release")
}
