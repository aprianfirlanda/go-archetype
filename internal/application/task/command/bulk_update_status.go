package taskcmd

import (
	"go-archetype/internal/domain/task"
)

type BulkUpdateStatus struct {
	PublicIDs []string
	Status    task.Status
}
