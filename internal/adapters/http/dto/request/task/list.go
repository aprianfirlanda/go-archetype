package taskreq

import (
	"errors"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
	"strings"
)

type List struct {
	Search   string `query:"search" validate:"omitempty,min=1,max=100"`
	Status   string `query:"status" validate:"omitempty,oneof=todo in_progress done"`
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Priority *int   `query:"priority" validate:"omitempty,min=1,max=5"`
}

func (q *List) Normalize() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}

	q.Search = strings.TrimSpace(q.Search)
}

func (q *List) ToListFilter() (taskquery.ListFilter, error) {
	var status task.Status

	if q.Status != "" {
		status = task.Status(q.Status)
		if !status.IsValid() {
			return taskquery.ListFilter{}, errors.New("invalid status")
		}
	}

	return taskquery.ListFilter{
		Search:   q.Search,
		Status:   status,
		Priority: q.Priority,
		Page:     q.Page,
		Limit:    q.Limit,
	}, nil
}
