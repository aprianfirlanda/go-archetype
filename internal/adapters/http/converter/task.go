package converter

import (
	"errors"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/domain/task"
)

func ToListFilter(q request.ListTasks) (task.ListFilter, error) {
	var status task.Status

	if q.Status != "" {
		status = task.Status(q.Status)
		if !status.IsValid() {
			return task.ListFilter{}, errors.New("invalid status")
		}
	}

	return task.ListFilter{
		Search:   q.Search,
		Status:   status,
		Priority: q.Priority,
		Page:     q.Page,
		Limit:    q.Limit,
	}, nil
}
