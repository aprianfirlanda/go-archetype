package tasksvc

import (
	"context"
	"go-archetype/internal/application/task/command"
	taskresult "go-archetype/internal/application/task/result"
	"go-archetype/internal/pkg/apperror"
)

func (s *service) BulkDelete(ctx context.Context, cmd taskcmd.BulkDelete) (*taskresult.BulkDeleteResult, error) {
	if len(cmd.PublicIDs) == 0 {
		return nil, apperror.Validation("no task publicIDs provided", nil)
	}

	result := &taskresult.BulkDeleteResult{
		Deleted: []string{},
		Failed:  []taskresult.BulkDeleteFailure{},
	}

	for _, publicID := range cmd.PublicIDs {

		if err := s.taskRepository.DeleteByPublicID(ctx, publicID); err != nil {
			result.Failed = append(result.Failed, taskresult.BulkDeleteFailure{
				PublicID: publicID,
				Reason:   err.Error(),
			})
			continue
		}

		result.Deleted = append(result.Deleted, publicID)
	}

	return result, nil
}
