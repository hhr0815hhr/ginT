package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hhr0815hhr/gint/internal"
	"github.com/hhr0815hhr/gint/internal/cache"
	"github.com/hhr0815hhr/gint/internal/config"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/hhr0815hhr/gint/internal/pkg/i18n"
	"github.com/hhr0815hhr/gint/internal/queue/memory_queue"
	"github.com/hhr0815hhr/gint/internal/queue/redis_queue"
)

func doInit() {
	i18n.InitI18n()
	internal.App = internal.InitApp()
	internal.App.Data["cache"] = cache.InitializeCache()
	initQueue()
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

func start() {
	app := internal.App
	port := config.Conf.Server.Port
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.Engine,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		<-quit
		log.Logger.Println("Shutting down server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Logger.Fatalf("Server forced to shutdown: %v", err)
		}
		log.Logger.Println("Server exiting")
	}()
	log.Logger.Printf("Starting server on :%d", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Logger.Fatalf(err.Error())
	}
}
