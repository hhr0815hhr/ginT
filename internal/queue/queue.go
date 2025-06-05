package queue

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hhr0815hhr/gint/internal"
)

// Message 定义队列中消息的结构
type Message struct {
	Body      gin.H
	MsgType   string //消息类型
	ReInCount int    // 当前重试次数  大于一定次数放入死信队列
	// 可以根据需要添加更丰富的元数据
	Headers map[string]interface{}
}

// Driver 队列驱动接口
type Driver interface {
	Publish(ctx context.Context, queueName string, message *Message) error
	Consume(ctx context.Context, queueName string, handler func(ctx context.Context, message *Message) error) error
	// 可选：添加队列管理方法，例如创建队列、删除队列等
	// EnsureQueueExists(ctx context.Context, queueName string) error
	// DeleteQueue(ctx context.Context, queueName string) error
	Close() error // 关闭连接，释放资源
}

const (
	DefaultQueueName = "default"
	DeadQueueName    = "dead_queue"
)

func PushQueue(msg *Message, queueName string) error {
	err := internal.App.Data["queue"].(Driver).Publish(context.Background(), queueName, msg)
	if err != nil {
		return err
	}
	return nil
}
