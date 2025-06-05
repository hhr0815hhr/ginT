// queue/redisqueue/redisqueue.go
package redis_queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hhr0815hhr/gint/internal/log"

	"github.com/go-redis/redis/v8"
	"github.com/hhr0815hhr/gint/internal/queue" // 替换为你的模块路径
)

// RedisListDriver 使用 Redis List 实现的队列驱动
type RedisListDriver struct {
	client redis.Cmdable
	Type   string
}

// NewRedisListDriver 创建一个新的 Redis List 驱动
func NewRedisListDriver(client redis.Cmdable, t string) *RedisListDriver {
	return &RedisListDriver{client: client, Type: t}
}

// Publish 将消息发布到 Redis List
func (d *RedisListDriver) Publish(ctx context.Context, queueName string, message *queue.Message) error {
	payloadJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}
	_, err = d.client.LPush(ctx, queueName, payloadJSON).Result()
	return err
}

// Consume 从 Redis List 消费消息
func (d *RedisListDriver) Consume(ctx context.Context, queueName string, handler func(ctx context.Context, message *queue.Message) error) error {
	for {
		result, err := d.client.BLPop(ctx, 0*time.Second, queueName).Result()
		if err != nil {
			if err != redis.Nil && err != context.Canceled {
				log.Logger.Printf("Error popping from Redis: %v\n", err)
				time.Sleep(time.Second)
			}
			if err == context.Canceled {
				return nil // 正常退出
			}
			continue
		}
		log.Logger.Printf("Received message from Redis: %v\n", result)
		if len(result) > 1 {
			payloadJSON := result[1]
			var payload = &queue.Message{}
			if err = json.Unmarshal([]byte(payloadJSON), payload); err != nil {
				log.Logger.Printf("Failed to unmarshal message payload: %v\n", err)
				continue
			}
			go handler(ctx, payload)
		}
	}
}

// Close 关闭 Redis 连接
func (d *RedisListDriver) Close() error {
	if d.Type == "cluster" {
		return d.client.(*redis.ClusterClient).Close()
	}
	return d.client.(*redis.Client).Close()
}
