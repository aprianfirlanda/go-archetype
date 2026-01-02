package tasksvc

import (
	"context"
	"errors"
	"go-archetype/internal/application/task/command"
	taskresult "go-archetype/internal/application/task/result"
)

func (s *Service) BulkDelete(ctx context.Context, cmd taskcmd.BulkDelete) (*taskresult.BulkDeleteResult, error) {
	if len(cmd.PublicIDs) == 0 {
		return nil, errors.New("no task publicIDs provided")
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
