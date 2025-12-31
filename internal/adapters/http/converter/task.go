package converter

import (
	"errors"
	"go-archetype/internal/adapters/http/dto/request"
	taskapp "go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
)

func ToListFilter(q request.ListTasks) (taskapp.ListFilter, error) {
	var status task.Status

	if q.Status != "" {
		status = task.Status(q.Status)
		if !status.IsValid() {
			return taskapp.ListFilter{}, errors.New("invalid status")
		}
	}

	return taskapp.ListFilter{
		Search:   q.Search,
		Status:   status,
		Priority: q.Priority,
		Page:     q.Page,
		Limit:    q.Limit,
	}, nil
}
