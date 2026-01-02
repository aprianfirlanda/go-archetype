package taskgorm

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"

	"gorm.io/gorm"
)

func (r *repository) FindByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	var model Task
	err := r.db.WithContext(ctx).First(&model, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, task.ErrNotFound
		}
		return nil, err
	}
	return toEntity(&model), nil
}
