package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) GetByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "GetByPublicID",
		"task_public_id": publicID,
	})

	if publicID == "" {
		log.Warn("task publicID is required")
		return nil, apperror.Validation("task publicID is required", nil)
	}

	entity, err := s.taskRepository.FindByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			log.WithError(err).Warn("task not found")
			return nil, apperror.NotFound("task not found", err)
		}
		log.WithError(err).Error("failed to get task from repository")
		return nil, apperror.Internal("failed to get task", err)
	}

	log.Info("task fetched")
	return entity, nil
}
