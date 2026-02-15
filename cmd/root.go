/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
		rel.FetchHelmRelease("", "")
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
}
