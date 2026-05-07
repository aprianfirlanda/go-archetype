package taskgorm

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"

	"gorm.io/gorm"
)

func (r *repository) FindByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "FindByPublicID",
		"task_public_id": publicID,
	})

	var model Model
	err := r.db.WithContext(ctx).First(&model, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("task record not found")
			return nil, task.ErrNotFound
		}
		log.WithError(err).Error("failed to query task record")
		return nil, err
	}
	log.Info("task record fetched")
	return toEntity(&model), nil
}
