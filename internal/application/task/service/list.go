package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) List(ctx context.Context, filter taskquery.ListFilter) ([]*task.Entity, int64, error) {
	tasks, total, err := s.taskRepository.FindAll(ctx, filter)
	if err != nil {
		return nil, 0, apperror.Internal("failed to list tasks", err)
	}
	return tasks, total, nil
}
