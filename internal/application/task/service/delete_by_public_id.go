package tasksvc

import (
	"context"
)

func (s *Service) DeleteByPublicID(ctx context.Context, publicID string) error {
	return s.taskRepository.DeleteByPublicID(ctx, publicID)
}
