package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) Create(ctx context.Context, cmd taskcmd.Create) (string, error) {
	log := componentLog(ctx).WithField("operation", "Create")

	entity := task.New(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)

	if err := entity.Validate(); err != nil {
		log.WithError(err).Warn("task validation failed")
		return "", apperror.Validation(err.Error(), err)
	}

	if err := s.taskRepository.Create(ctx, entity); err != nil {
		log.WithError(err).Error("failed to create task in repository")
		return "", apperror.Internal("failed to create task", err)
	}

	log.WithField("task_public_id", entity.PublicID).Info("task created")
	return entity.PublicID, nil
}
