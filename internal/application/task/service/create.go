package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) Create(ctx context.Context, cmd taskcmd.Create) (string, error) {
	entity := task.New(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)

	if err := entity.Validate(); err != nil {
		return "", apperror.Validation(err.Error(), err)
	}

	if err := s.taskRepository.Create(ctx, entity); err != nil {
		return "", apperror.Internal("failed to create task", err)
	}

	return entity.PublicID, nil
}
