package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) UpdateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	taskEntity, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			return apperror.NotFound("task not found", err)
		}
		return apperror.Internal("failed to get task", err)
	}

	if err := taskEntity.UpdateStatus(cmd.Status); err != nil {
		return apperror.Validation(err.Error(), err)
	}

	if err := s.taskRepository.UpdateByPublicID(ctx, taskEntity); err != nil {
		return apperror.Internal("failed to update task status", err)
	}

	return nil
}
