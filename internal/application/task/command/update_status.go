package taskcmd

import (
	"go-archetype/internal/domain/task"
)

type UpdateStatus struct {
	PublicID string
	Status   task.Status
}
