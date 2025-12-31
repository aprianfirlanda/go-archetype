package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
)

func (s *Service) GetByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	if publicID == "" {
		return nil, errors.New("task publicID is required")
	}
	return s.taskRepository.FindByPublicID(ctx, publicID)
}
