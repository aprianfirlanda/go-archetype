package portin

import "context"

type MessageHandler func(ctx context.Context, payload []byte) error

type MessageConsumer interface {
	Consume(ctx context.Context, topic string, handler MessageHandler) error
}
