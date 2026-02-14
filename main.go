package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "helm-health",
		Short: "Helm plugin for health checking",
		Long:  `A Helm plugin that provides health checking capabilities for deployed applications.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("helm-health plugin initialized")
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
