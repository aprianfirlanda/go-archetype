package portout

import (
	"context"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
)

type TaskRepository interface {
	Create(ctx context.Context, task *task.Entity) error
	FindByPublicID(ctx context.Context, publicID string) (*task.Entity, error)
	FindAll(ctx context.Context, filter taskquery.ListFilter) ([]*task.Entity, int64, error)
	UpdateByPublicID(ctx context.Context, task *task.Entity) error
	UpdateStatusByPublicID(ctx context.Context, publicID string, status task.Status) error
	DeleteByPublicID(ctx context.Context, publicID string) error
	BulkUpdateStatus(ctx context.Context, publicIDs []string, status task.Status) error
	BulkDelete(ctx context.Context, publicIDs []string) error
}
