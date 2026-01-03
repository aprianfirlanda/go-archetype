package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *Service) DeleteByPublicID(ctx context.Context, publicID string) error {
	if publicID == "" {
		return apperror.Validation("task publicID is required", nil)
	}

	if err := s.taskRepository.DeleteByPublicID(ctx, publicID); err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return apperror.NotFound("task not found", err)
		}
		return apperror.Internal("failed to delete task", err)
	}

	return nil
}
