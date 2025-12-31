package tasksvc

import (
	"context"
	"errors"
	taskcmd "go-archetype/internal/application/task/command"
)

func (s *Service) BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) error {
	if len(cmd.PublicIDs) == 0 {
		return errors.New("no task publicIDs provided")
	}

	if !cmd.Status.IsValid() {
		return errors.New("invalid status")
	}

	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update each task with business validation
	for _, publicID := range cmd.PublicIDs {
		if err := s.updateStatus(ctx, taskcmd.UpdateStatus{
			PublicID: publicID,
			Status:   cmd.Status,
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}
