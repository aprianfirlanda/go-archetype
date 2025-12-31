package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
)

func (s *Service) List(ctx context.Context, filter taskquery.ListFilter) ([]*task.Entity, int64, error) {
	return s.taskRepository.FindAll(ctx, filter)
}
