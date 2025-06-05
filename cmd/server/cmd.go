package server

import (
	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Gin HTTP server",
	Long:  `Starts the Gin HTTP server on the specified port.`,
	Run: func(cmd *cobra.Command, args []string) {
		doInit()
		//http.Serve()
		//runGlobalGoroutines()
		start()
	},
}
