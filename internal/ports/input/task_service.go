package input

import (
	"context"
	"go-archetype/internal/domain/task"
)

type TaskService interface {
	CreateTask(ctx context.Context, task *task.Entity) (string, error)
	GetTaskByPublicID(ctx context.Context, id string) (*task.Entity, error)
	ListTasks(ctx context.Context, filter task.ListFilter) ([]*task.Entity, int64, error)
	UpdateTask(ctx context.Context, task *task.Entity) error
	UpdateTaskStatus(ctx context.Context, id string, status task.Status) error
	DeleteTaskByPublicID(ctx context.Context, id string) error
	BulkUpdateStatus(ctx context.Context, ids []string, status task.Status) error
	BulkDelete(ctx context.Context, ids []string) error
}
