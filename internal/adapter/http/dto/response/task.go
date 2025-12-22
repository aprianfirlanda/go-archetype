package response

import (
	"go-archetype/internal/domain/task"
	"time"
)

type CreateTaskValidateError struct {
	Title       []string `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Priority    []string `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type Task struct {
	ID          string      `json:"id,omitempty"`
	Title       string      `json:"title,omitempty"`
	Description string      `json:"description,omitempty"`
	Status      task.Status `json:"status,omitempty"`
	Priority    int         `json:"priority,omitempty"`
	DueDate     *time.Time  `json:"due_date,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Completed   bool        `json:"completed,omitempty"`
	CreatedAt   time.Time   `json:"created_at,omitempty"`
}

type ListTasksValidateError struct {
	Search   string   `query:"search,omitempty"`
	Status   string   `query:"status,omitempty"`
	Page     []string `query:"page,omitempty"`
	Limit    []string `query:"limit,omitempty"`
	Priority []string `query:"priority,omitempty"`
}

type UpdateTaskValidateError struct {
	Title       []string `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Priority    []string `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdateTaskStatusValidateError struct {
	Status []string `json:"status,omitempty"`
}

type BulkUpdateTaskStatusValidateError struct {
	IDs    []string `json:"ids,omitempty"`
	Status []string `json:"status,omitempty"`
}

type BulkDeleteTasksValidateError struct {
	IDs []string `json:"ids,omitempty"`
}
