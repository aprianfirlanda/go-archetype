package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	taskresult "go-archetype/internal/application/task/result"
)

func (s *Service) BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) (*taskresult.BulkUpdateStatusResult, error) {
	if len(cmd.PublicIDs) == 0 {
		return nil, errors.New("no task publicIDs provided")
	}

	if !cmd.Status.IsValid() {
		return nil, errors.New("invalid status")
	}

	result := &taskresult.BulkUpdateStatusResult{
		Updated: []string{},
		Failed:  []taskresult.BulkUpdateStatusFail{},
	}

	for _, publicID := range cmd.PublicIDs {
		task, err := s.taskRepository.FindByPublicID(ctx, publicID)
		if err != nil {
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		if err := task.UpdateStatus(cmd.Status); err != nil {
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		if err := s.taskRepository.UpdateByPublicID(ctx, task); err != nil {
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		result.Updated = append(result.Updated, publicID)
	}

	return result, nil
}
