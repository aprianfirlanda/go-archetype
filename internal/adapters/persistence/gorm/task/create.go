package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) Create(ctx context.Context, t *task.Entity) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "Create",
		"task_public_id": t.PublicID,
	})

	model := toModel(t)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		log.WithError(err).Error("failed to create task record")
		return err
	}

	log.Info("task record created")
	return nil
}
