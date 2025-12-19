package request

import (
	"go-archetype/internal/domain/task"
	"time"
)

type CreateTask struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description"`
	Priority    int        `json:"priority" validate:"min=1,max=5"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags"`
}

type UpdateTask struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *int       `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
	Tags        *[]string  `json:"tags"`
}

type ListTasks struct {
	Search   string `query:"search"`
	Status   string `query:"status"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
	Priority *int   `query:"priority"`
}

type UpdateTaskStatus struct {
	Status task.Status `json:"status" validate:"required,oneof=todo in_progress done"`
}

type BulkUpdateTaskStatus struct {
	IDs    []string    `json:"ids" validate:"required,min=1,dive,required"`
	Status task.Status `json:"status" validate:"required,oneof=todo in_progress done"`
}

type BulkDeleteTasks struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,required"`
}
