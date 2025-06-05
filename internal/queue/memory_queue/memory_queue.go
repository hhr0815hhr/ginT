package memory_queue

import (
	"context"
	"sync"

	"github.com/hhr0815hhr/gint/internal/queue" // 替换为你的模块路径
)

// InMemoryDriver 使用内存 Channel 实现的队列驱动
type InMemoryDriver struct {
	queues map[string]chan *queue.Message
	mu     sync.Mutex
}

// NewInMemoryDriver 创建一个新的内存队列驱动
func NewInMemoryDriver() *InMemoryDriver {
	return &InMemoryDriver{
		queues: make(map[string]chan *queue.Message),
	}
}

// Publish 将消息发布到内存队列
func (d *InMemoryDriver) Publish(ctx context.Context, queueName string, message *queue.Message) error {
	d.mu.Lock()
	q, ok := d.queues[queueName]
	if !ok {
		q = make(chan *queue.Message, 100) // 可配置缓冲区大小
		d.queues[queueName] = q
	}
	d.mu.Unlock()

	select {
	case q <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Consume 从内存队列消费消息
func (d *InMemoryDriver) Consume(ctx context.Context, queueName string, handler func(ctx context.Context, message *queue.Message) error) error {
	d.mu.Lock()
	q, ok := d.queues[queueName]
	if !ok {
		q = make(chan *queue.Message, 100) // 确保队列存在
		d.queues[queueName] = q
	}
	d.mu.Unlock()

	for {
		select {
		case msg := <-q:
			go handler(ctx, msg)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Close 关闭内存队列 (对于内存队列，通常不需要显式关闭)
func (d *InMemoryDriver) Close() error {
	return nil
}
