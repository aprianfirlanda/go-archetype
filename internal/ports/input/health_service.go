package portin

import (
	"context"
)

type HealthService interface {
	Liveness(ctx context.Context) bool
	Readiness(ctx context.Context) error
}
