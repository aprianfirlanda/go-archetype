package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) Update(ctx context.Context, cmd taskcmd.Update) error {
	taskEntity, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return apperror.NotFound("task not found", err)
		}
		return apperror.Internal("failed to get task", err)
	}

	taskEntity.Update(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)
	if err := s.taskRepository.UpdateByPublicID(ctx, taskEntity); err != nil {
		return apperror.Internal("failed to update task", err)
	}

	return nil
}
