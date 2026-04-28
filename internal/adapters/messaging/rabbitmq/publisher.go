package messagingrmq

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	ch *amqp091.Channel
}

func NewPublisher(conn *Connection) (*Publisher, error) {
	ch, err := conn.Conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Publisher{ch: ch}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	return p.ch.PublishWithContext(
		ctx,
		"",    // default exchange
		topic, // routing key = queue name
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        payload,
		},
	)
}
