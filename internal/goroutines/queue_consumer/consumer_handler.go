package queue_consumer

import (
	"github.com/hhr0815hhr/gint/internal/queue"
)

type ConsumerHandler struct {
}

func commonRetry(message *queue.Message) error {
	message.ReInCount++
	return queue.PushQueue(message, queue.DefaultQueueName)
}

func (c *ConsumerHandler) handleTest(message *queue.Message) error {
	return nil
}
