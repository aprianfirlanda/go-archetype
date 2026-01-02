package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
)

func (s *Service) Create(ctx context.Context, cmd taskcmd.Create) (string, error) {
	entity := task.New(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)

	if err := entity.Validate(); err != nil {
		return "", err
	}

	if err := s.taskRepository.Create(ctx, entity); err != nil {
		return "", err
	}

	return entity.PublicID, nil
}
