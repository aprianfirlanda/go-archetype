package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) DeleteByPublicID(ctx context.Context, publicID string) error {
	result := r.db.WithContext(ctx).
		Where("public_id = ?", publicID).
		Delete(&Model{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return task.ErrNotFound
	}

	return nil
}
