package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) GetByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	if publicID == "" {
		return nil, apperror.Validation("task publicID is required", nil)
	}

	entity, err := s.taskRepository.FindByPublicID(ctx, publicID)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return nil, apperror.NotFound("task not found", err)
		}
		return nil, apperror.Internal("failed to get task", err)
	}

	return entity, nil
}
