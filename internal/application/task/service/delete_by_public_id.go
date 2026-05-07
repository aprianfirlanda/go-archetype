package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) DeleteByPublicID(ctx context.Context, publicID string) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "DeleteByPublicID",
		"task_public_id": publicID,
	})

	if publicID == "" {
		log.Warn("task publicID is required")
		return apperror.Validation("task publicID is required", nil)
	}

	if err := s.taskRepository.DeleteByPublicID(ctx, publicID); err != nil {
		if errors.Is(err, task.ErrNotFound) {
			log.WithError(err).Warn("task not found")
			return apperror.NotFound("task not found", err)
		}
		log.WithError(err).Error("failed to delete task in repository")
		return apperror.Internal("failed to delete task", err)
	}

	log.Info("task deleted")
	return nil
}
