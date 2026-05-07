package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) DeleteByPublicID(ctx context.Context, publicID string) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "DeleteByPublicID",
		"task_public_id": publicID,
	})

	result := r.db.WithContext(ctx).
		Where("public_id = ?", publicID).
		Delete(&Model{})

	if result.Error != nil {
		log.WithError(result.Error).Error("failed to delete task record")
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Warn("task record not found")
		return task.ErrNotFound
	}

	log.Info("task record deleted")
	return nil
}
