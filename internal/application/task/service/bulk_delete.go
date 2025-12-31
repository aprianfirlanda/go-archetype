package tasksvc

import (
	"context"
	"errors"
	taskcmd "go-archetype/internal/application/task/command"
)

func (s *Service) BulkDelete(ctx context.Context, cmd taskcmd.BulkDelete) error {
	if len(cmd.PublicIDs) == 0 {
		return errors.New("no task publicIDs provided")
	}

	return s.taskRepository.BulkDelete(ctx, cmd.PublicIDs)
}
