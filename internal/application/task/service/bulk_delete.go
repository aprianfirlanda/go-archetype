package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/application/task/result"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) BulkDelete(ctx context.Context, cmd taskcmd.BulkDelete) (*taskresult.BulkDeleteResult, error) {
	log := componentLog(ctx).WithField("operation", "BulkDelete")

	if len(cmd.PublicIDs) == 0 {
		log.Warn("no task publicIDs provided")
		return nil, apperror.Validation("no task publicIDs provided", nil)
	}

	result := &taskresult.BulkDeleteResult{
		Deleted: []string{},
		Failed:  []taskresult.BulkDeleteFailure{},
	}

	for _, publicID := range cmd.PublicIDs {
		if err := s.taskRepository.DeleteByPublicID(ctx, publicID); err != nil {
			log.WithError(err).WithField("task_public_id", publicID).Warn("failed to delete task")
			result.Failed = append(result.Failed, taskresult.BulkDeleteFailure{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		result.Deleted = append(result.Deleted, publicID)
	}

	log.WithFields(map[string]any{
		"deleted_count": len(result.Deleted),
		"failed_count":  len(result.Failed),
	}).Info("bulk delete completed")

	return result, nil
}
