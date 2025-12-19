package response

import (
	"go-archetype/internal/domain/task"
	"time"
)

type Task struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Status      task.Status `json:"status"`
	Priority    int         `json:"priority"`
	DueDate     *time.Time  `json:"due_date"`
	Tags        []string    `json:"tags"`
	Completed   bool        `json:"completed"`
	CreatedAt   time.Time   `json:"created_at"`
}
