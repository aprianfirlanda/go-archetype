package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/application/task/result"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) BulkUpdateStatus(ctx context.Context, cmd taskcmd.BulkUpdateStatus) (*taskresult.BulkUpdateStatusResult, error) {
	log := componentLog(ctx).WithFields(map[string]any{
		"operation": "BulkUpdateStatus",
		"status":    cmd.Status,
	})

	if len(cmd.PublicIDs) == 0 {
		log.Warn("no task publicIDs provided")
		return nil, apperror.Validation("no task publicIDs provided", nil)
	}

	if !cmd.Status.IsValid() {
		log.Warn("invalid status")
		return nil, apperror.Validation("invalid status", nil)
	}

	result := &taskresult.BulkUpdateStatusResult{
		Updated: []string{},
		Failed:  []taskresult.BulkUpdateStatusFail{},
	}

	for _, publicID := range cmd.PublicIDs {
		task, err := s.taskRepository.FindByPublicID(ctx, publicID)
		if err != nil {
			log.WithError(err).WithField("task_public_id", publicID).Warn("failed to get task")
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		if err := task.UpdateStatus(cmd.Status); err != nil {
			log.WithError(err).WithField("task_public_id", publicID).Warn("failed to update task status")
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		if err := s.taskRepository.UpdateByPublicID(ctx, task); err != nil {
			log.WithError(err).WithField("task_public_id", publicID).Warn("failed to persist task status")
			result.Failed = append(result.Failed, taskresult.BulkUpdateStatusFail{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		result.Updated = append(result.Updated, publicID)
	}

	log.WithFields(map[string]any{
		"updated_count": len(result.Updated),
		"failed_count":  len(result.Failed),
	}).Info("bulk update status completed")

	return result, nil
}
