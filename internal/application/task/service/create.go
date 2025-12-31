package tasksvc

import (
	"context"
	taskcmd "go-archetype/internal/application/task/command"
	taskDomain "go-archetype/internal/domain/task"
)

func (s *Service) Create(ctx context.Context, cmd taskcmd.Create) (string, error) {
	task := taskDomain.New(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)

	if err := task.Validate(); err != nil {
		return "", err
	}

	if err := s.taskRepository.Create(ctx, task); err != nil {
		return "", err
	}

	return task.PublicID, nil
}
