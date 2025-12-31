package request

import (
	"strings"
	"time"
)

type CreateTask struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description" validate:"max=500"`
	Priority    int        `json:"priority" validate:"min=1,max=5"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags" validate:"max=10,dive,min=1"`
}

type ListTasks struct {
	Search   string `query:"search" validate:"omitempty,min=1,max=100"`
	Status   string `query:"status" validate:"omitempty,oneof=todo in_progress done"`
	Page     int    `query:"page" validate:"omitempty,min=1"`
	Limit    int    `query:"limit" validate:"omitempty,min=1,max=100"`
	Priority *int   `query:"priority" validate:"omitempty,min=1,max=5"`
}

func (q *ListTasks) Normalize() {
	if q.Page == 0 {
		q.Page = 1
	}
	if q.Limit == 0 {
		q.Limit = 10
	}

	q.Search = strings.TrimSpace(q.Search)
}

type UpdateTask struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description" validate:"max=500"`
	Priority    int        `json:"priority" validate:"required,min=1,max=5"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags" validate:"max=10,dive,min=1"`
}

type UpdateTaskStatus struct {
	Status string `json:"status" validate:"required,oneof=todo in_progress done"`
}

type BulkUpdateTaskStatus struct {
	IDs    []string `json:"ids" validate:"required,min=1,dive,required"`
	Status string   `json:"status" validate:"required,oneof=todo in_progress done"`
}

type BulkDeleteTasks struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,required"`
}
