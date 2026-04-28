package messagingrmq

import (
	"context"
	"go-archetype/internal/ports/input"
	"log"

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
			err := handler(ctx, msg.Body)
			if err != nil {
				log.Println("handler error:", err)
				_ = msg.Nack(false, true) // retry
				continue
			}
			_ = msg.Ack(false)
		}
	}()

	return nil
}
