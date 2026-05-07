package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) List(ctx context.Context, filter taskquery.ListFilter) ([]*task.Entity, int64, error) {
	log := componentLog(ctx).WithField("operation", "List")

	tasks, total, err := s.taskRepository.FindAll(ctx, filter)
	if err != nil {
		log.WithError(err).Error("failed to list tasks from repository")
		return nil, 0, apperror.Internal("failed to list tasks", err)
	}

	log.WithFields(map[string]any{
		"count": total,
		"page":  filter.Page,
		"limit": filter.Limit,
	}).Info("tasks listed")

	return tasks, total, nil
}
