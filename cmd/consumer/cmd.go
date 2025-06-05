package consumer

import (
	"github.com/hhr0815hhr/gint/internal"
	"github.com/hhr0815hhr/gint/internal/cache"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/goroutines"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/hhr0815hhr/gint/internal/pkg/i18n"
	"github.com/hhr0815hhr/gint/internal/queue/memory_queue"
	"github.com/hhr0815hhr/gint/internal/queue/redis_queue"
	"github.com/spf13/cobra"
)

var ConsumeCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Start consumers",
	Long:  `Starts every topic's consumer`,
	Run: func(cmd *cobra.Command, args []string) {
		doInit()
		startConsumer()
	},
}

func initQueue() {
	switch config.Conf.Server.Queue {
	case "memory":
		internal.App.Data["queue"] = memory_queue.NewInMemoryDriver()
	case "redis":
		internal.App.Data["queue"] = redis_queue.NewRedisListDriver(cache.Client, config.Conf.Redis.Type)
	default:
		log.Logger.Fatalf("unknown queue driver: %s", config.Conf.Server.Queue)
	}
	log.Logger.Println("初始化队列...success")
}
func doInit() {
	i18n.InitI18n()
	internal.App = internal.InitApp()
	internal.App.Data["cache"] = cache.InitializeCache()
	initQueue()
}

func startConsumer() {
	goroutines.RunGlobalGoroutines()
}
