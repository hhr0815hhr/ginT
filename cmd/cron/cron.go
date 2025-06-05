package cron

import (
	"github.com/hhr0815hhr/gint/internal"
	"github.com/hhr0815hhr/gint/internal/cache"
	cron2 "github.com/hhr0815hhr/gint/internal/cron"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/spf13/cobra"
)

var CronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Start cron jobs",
	Long:  `Start cron jobs`,
	Run: func(cmd *cobra.Command, args []string) {
		//初始化rcon配置
		doInit()
		startCronJob()
	},
}

func doInit() {
	internal.App = internal.InitApp()
	internal.App.Data["cache"] = cache.InitializeCache()
}

func startCronJob() {
	log.Logger.Info("starting cron jobs...")

	c := cron2.New()
	c.AddFuncs()
	c.Start()
	select {}
}
