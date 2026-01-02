package taskgorm

import (
	"context"
	"go-archetype/internal/domain/task"
)

func (r *repository) UpdateByPublicID(ctx context.Context, t *task.Entity) error {
	result := r.db.WithContext(ctx).
		Model(&Task{}).
		Where("public_id = ?", t.PublicID).
		Updates(toModel(t))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return task.ErrNotFound
	}

	return nil
}
