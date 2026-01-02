package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
	"time"
)

func (r *repository) UpdateStatusByPublicID(ctx context.Context, publicID string, status task.Status) error {
	result := r.db.WithContext(ctx).
		Model(&Model{}).
		Where("public_id = ?", publicID).
		Updates(map[string]interface{}{
			"status":     string(status),
			"completed":  status == task.StatusDone,
			"updated_at": time.Now(),
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return task.ErrNotFound
	}

	return nil
}
