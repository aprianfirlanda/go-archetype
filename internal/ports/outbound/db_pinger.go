package outbound

import "context"

type DBPinger interface {
	Ping(ctx context.Context) error
}
