package goroutines

import (
	"github.com/hhr0815hhr/gint/internal"
	"github.com/hhr0815hhr/gint/internal/goroutines/queue_consumer"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/hhr0815hhr/gint/internal/queue"
)

func RunGlobalGoroutines() {
	log.Logger.Println("[goroutine]消息队列消费者启动...success")
	// 优化为遍历topic切片，每个topic启动n个消费者
	queue_consumer.Consumer(queue.DefaultQueueName, internal.App.Data["queue"].(queue.Driver))
}
