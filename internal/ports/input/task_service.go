package portin

import (
	"context"
	taskcmd "go-archetype/internal/application/task/command"
	taskquery "go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
)

type TaskService interface {
	Create(ctx context.Context, cmd taskcmd.Create) (string, error)
	GetByPublicID(ctx context.Context, id string) (*task.Entity, error)
	List(ctx context.Context, query taskquery.ListFilter) ([]*task.Entity, int64, error)
	Update(ctx context.Context, cmd taskcmd.Update) error
	UpdateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error
	BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) error
	DeleteByPublicID(ctx context.Context, id string) error
	BulkDelete(ctx context.Context, cmd taskcmd.BulkDelete) error
}
