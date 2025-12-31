package tasksvc

import (
	"context"
	taskcmd "go-archetype/internal/application/task/command"
)

func (s *Service) Update(ctx context.Context, cmd taskcmd.Update) error {
	task, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		return err
	}

	task.Update(
		cmd.Title,
		cmd.Description,
		cmd.Priority,
		cmd.DueDate,
		cmd.Tags,
	)
	return s.taskRepository.UpdateByPublicID(ctx, task)
}
