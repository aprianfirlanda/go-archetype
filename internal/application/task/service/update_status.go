package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) UpdateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation":      "UpdateStatus",
		"task_public_id": cmd.PublicID,
		"status":         cmd.Status,
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

	if err := taskEntity.UpdateStatus(cmd.Status); err != nil {
		log.WithError(err).Warn("invalid task status update")
		return apperror.Validation(err.Error(), err)
	}

	if err := s.taskRepository.UpdateByPublicID(ctx, taskEntity); err != nil {
		log.WithError(err).Error("failed to update task status in repository")
		return apperror.Internal("failed to update task status", err)
	}

	log.Info("task status updated")
	return nil
}
