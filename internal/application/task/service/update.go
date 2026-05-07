package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) Update(ctx context.Context, cmd taskcmd.Update) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "Update",
		"task_public_id": cmd.PublicID,
	})

	taskEntity, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		if errors.Is(err, task.ErrNotFound) {
			log.WithError(err).Warn("task not found")
			return apperror.NotFound("task not found", err)
		}
		log.WithError(err).Error("failed to get task from repository")
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
		log.WithError(err).Error("failed to update task in repository")
		return apperror.Internal("failed to update task", err)
	}

	log.Info("task updated")
	return nil
}
