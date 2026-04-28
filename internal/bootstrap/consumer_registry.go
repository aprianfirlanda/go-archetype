package bootstrap

import (
	"context"
	portin "go-archetype/internal/ports/input"
)

type ConsumerRegistration struct {
	Topic   string
	Handler portin.MessageHandler
}

type ConsumerRegistry struct {
	consumers []ConsumerRegistration
}

func NewConsumerRegistry() *ConsumerRegistry {
	return &ConsumerRegistry{}
}

func (r *ConsumerRegistry) Register(topic string, handler portin.MessageHandler) {
	r.consumers = append(r.consumers, ConsumerRegistration{
		Topic:   topic,
		Handler: handler,
	})
}

func (r *ConsumerRegistry) Start(ctx context.Context, consumer portin.MessageConsumer) error {
	for _, c := range r.consumers {
		if err := consumer.Consume(ctx, c.Topic, c.Handler); err != nil {
			return err
		}
	}
	return nil
}
