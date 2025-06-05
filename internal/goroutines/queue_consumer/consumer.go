package queue_consumer

import (
	"context"
	"fmt"

	_const "github.com/hhr0815hhr/gint/internal/const"
	"github.com/hhr0815hhr/gint/internal/log"
	"github.com/hhr0815hhr/gint/internal/queue"
)

var hd = &ConsumerHandler{}

func Consumer(queueName string, driver queue.Driver) {
	ctx := context.Background()
	handler := func(ctx context.Context, message *queue.Message) error {
		fmt.Printf("Consumed type: %s, data: %v, Headers: %v\n", message.MsgType, message.Body, message.Headers)
		if message.ReInCount > 3 {
			// 放入死信队列
			return queue.PushQueue(message, queue.DeadQueueName)
		}
		//todo 按事件类型处理
		switch message.MsgType {

		case _const.QUEUE_TEST:
			return hd.handleTest(message)

		default:
			log.Logger.Printf("Unknown message type: %s\n", message.MsgType)
		}
		return nil
	}
	err := driver.Consume(ctx, queueName, handler)
	if err != nil {
		log.Logger.Fatalf("Error consuming: %v\n", err)
	}
}
