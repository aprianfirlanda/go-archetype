package taskquery

import "go-archetype/internal/domain/task"

type ListFilter struct {
	Search   string
	Status   task.Status
	Priority *int
	Page     int
	Limit    int
}
