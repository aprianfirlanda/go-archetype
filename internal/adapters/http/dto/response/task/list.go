package taskresp

import (
	"go-archetype/internal/domain/task"
	"time"
)

type ListValidateError struct {
	Search   string   `query:"search,omitempty"`
	Status   string   `query:"status,omitempty"`
	Page     []string `query:"page,omitempty"`
	Limit    []string `query:"limit,omitempty"`
	Priority []string `query:"priority,omitempty"`
}

type ListItem struct {
	ID          string     `json:"id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Status      string     `json:"status,omitempty"`
	Priority    int        `json:"priority,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Tags        []string   `json:"tags,omitempty"`
	Completed   bool       `json:"completed,omitempty"`
	CreatedAt   time.Time  `json:"created_at,omitempty"`
}

func ToListItem(entity *task.Entity) ListItem {
	return ListItem{
		ID:          entity.PublicID,
		Title:       entity.Title,
		Description: entity.Description,
		Status:      entity.Status.String(),
		Priority:    entity.Priority,
		DueDate:     entity.DueDate,
		Tags:        entity.Tags,
		Completed:   entity.Completed,
		CreatedAt:   entity.CreatedAt,
	}
}

func ToList(items []*task.Entity) []ListItem {
	result := make([]ListItem, len(items))

	for i, entity := range items {
		result[i] = ToListItem(entity)
	}

	return result
}
