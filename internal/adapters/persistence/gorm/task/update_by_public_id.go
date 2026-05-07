package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) UpdateByPublicID(ctx context.Context, t *task.Entity) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "UpdateByPublicID",
		"task_public_id": t.PublicID,
	})

	result := r.db.WithContext(ctx).
		Model(&Model{}).
		Where("public_id = ?", t.PublicID).
		Updates(toModel(t))

	if result.Error != nil {
		log.WithError(result.Error).Error("failed to update task record")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Warn("task record not found")
		return task.ErrNotFound
	}

	log.Info("task record updated")
	return nil
}
