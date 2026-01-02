package taskresp

import (
	"go-archetype/internal/domain/task"
	"time"
)

type Detail struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    int        `json:"priority"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToDetail(entity *task.Entity) Detail {
	return Detail{
		ID:          entity.PublicID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      entity.Status.String(),
		Priority:    entity.Priority,
		DueDate:     entity.DueDate,
		Tags:        entity.Tags,
		Completed:   entity.Completed,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
