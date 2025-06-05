package main

import (
	"fmt"
	"os"

	"github.com/hhr0815hhr/gint/cmd/consumer"
	"github.com/hhr0815hhr/gint/cmd/cron"
	"github.com/hhr0815hhr/gint/cmd/gen"
	"github.com/hhr0815hhr/gint/cmd/server"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gint",
	Short: "gint服务",
	Long:  `gint服务`,
}

func init() {
	rootCmd.AddCommand(server.ServeCmd)
	rootCmd.AddCommand(gen.GenCmd)
	rootCmd.AddCommand(consumer.ConsumeCmd)
	rootCmd.AddCommand(cron.CronCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
