package portout

import "context"

type MessagePublisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}
