package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) Create(ctx context.Context, t *task.Entity) error {
	model := toModel(t)
	return r.db.WithContext(ctx).Create(model).Error
}
