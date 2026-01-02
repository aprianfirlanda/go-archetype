package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
)

func (s *Service) BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) error {
	if len(cmd.PublicIDs) == 0 {
		return errors.New("no task publicIDs provided")
	}

	if !cmd.Status.IsValid() {
		return errors.New("invalid status")
	}

	return s.taskRepository.BulkUpdateStatus(ctx, cmd.PublicIDs, cmd.Status)
}
