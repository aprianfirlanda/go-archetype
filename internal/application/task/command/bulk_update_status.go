package taskcmd

import (
	task2 "go-archetype/internal/domain/task"
)

type BulkUpdateStatus struct {
	PublicIDs []string
	Status    task2.Status
}
