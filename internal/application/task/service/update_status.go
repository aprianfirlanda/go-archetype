package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
)

func (s *Service) UpdateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	task, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		return err
	}

	if err := task.UpdateStatus(cmd.Status); err != nil {
		return err
	}

	return s.taskRepository.UpdateByPublicID(ctx, task)
}
