package tasksvc

import (
	"context"
	taskcmd "go-archetype/internal/application/task/command"
)

func (s *Service) updateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	task, err := s.taskRepository.FindByPublicID(ctx, cmd.PublicID)
	if err != nil {
		return err
	}

	if err := task.UpdateStatus(cmd.Status); err != nil {
		return err
	}

	return s.taskRepository.Update(ctx, task)
}

func (s *Service) UpdateStatusSingle(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	return s.updateStatus(ctx, cmd)
}
