package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	taskresult "go-archetype/internal/application/task/result"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) (*taskresult.BulkUpdateStatusResult, error) {
	if len(cmd.PublicIDs) == 0 {
		return nil, apperror.Validation("no task publicIDs provided", nil)
	}

	if !cmd.Status.IsValid() {
		return nil, apperror.Validation("invalid status", nil)
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
