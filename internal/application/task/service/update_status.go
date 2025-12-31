package tasksvc

import (
	"context"
	taskcmd "go-archetype/internal/application/task/command"
)

func (s *Service) UpdateStatus(ctx context.Context, cmd taskcmd.UpdateStatus) error {
	return s.taskRepository.UpdateStatusByPublicID(ctx, cmd.PublicID, cmd.Status)
}
