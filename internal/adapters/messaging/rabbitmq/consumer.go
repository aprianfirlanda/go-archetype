package messagingrmq

import (
	"context"
	"fmt"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"

	"github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	ch *amqp091.Channel
}

func NewConsumer(conn *Connection) (*Consumer, error) {
	ch, err := conn.Conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Consumer{ch: ch}, nil
}

func (c *Consumer) Consume(
	ctx context.Context,
	topic string,
	handler portin.MessageHandler,
) error {

	_, err := c.ch.QueueDeclare(
		topic,
		true, // durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.ch.Consume(
		topic,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			rid := msg.CorrelationId
			if rid == "" {
				rid = msg.MessageId
			}
			if rid == "" {
				rid = fmt.Sprintf("msg-%d", msg.DeliveryTag)
			}

			messageLog := logging.ComponentLogger(logging.FromContext(ctx), "messaging.rabbitmq.consumer").
				WithFields(map[string]any{
					"topic":      topic,
					"rid":        rid,
					"request_id": rid,
				})
			msgCtx := logging.WithLogger(ctx, messageLog)
			msgCtx = logging.WithRequestID(msgCtx, rid)

			err := handler(msgCtx, msg.Body)
			if err != nil {
				messageLog.WithError(err).Error("consumer handler error")
				_ = msg.Nack(false, true) // retry
				continue
			}

			messageLog.Info("message processed")
			_ = msg.Ack(false)
		}
	}()

	return nil
}
