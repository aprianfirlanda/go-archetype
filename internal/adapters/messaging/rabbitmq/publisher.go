package messagingrmq

import (
	"context"
	"go-archetype/internal/infrastructure/logging"

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
	log := logging.ComponentLogger(logging.FromContext(ctx), "messaging.rabbitmq.publisher").WithField("topic", topic)
	rid := logging.RequestIDFromContext(ctx)
	if err := p.ch.PublishWithContext(
		ctx,
		"",    // default exchange
		topic, // routing key = queue name
		false,
		false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          payload,
			CorrelationId: rid,
			MessageId:     rid,
		},
	); err != nil {
		log.WithError(err).Error("failed to publish message")
		return err
	}

	log.Info("message published")
	return nil
}
