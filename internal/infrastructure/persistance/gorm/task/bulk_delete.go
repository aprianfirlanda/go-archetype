package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) BulkDelete(ctx context.Context, publicIDs []string) error {
	result := r.db.WithContext(ctx).
		Where("public_id IN ?", publicIDs).
		Delete(&Task{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return task.ErrNotFound
	}

	return nil
}
