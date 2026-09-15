package cmd

import (
	"fmt"
	"os"

	"github.com/amankr1098/helm-health/internal/server"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the helm-health HTTP API server",
	Long: `Starts an HTTP server that exposes helm-health results as JSON.

Endpoints:
  GET /api/health?release=<name>&namespace=<ns>   Health report for one release
  GET /api/releases?namespace=<ns>                List releases (all namespaces if omitted)`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		if err := server.Start(":" + port); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("port", "p", "8000", "Port to listen on")
}
