package taskgorm

import (
	"go-archetype/internal/domain/task"
	"strings"
	"time"
)

// Model is the GORM model (infrastructure concern)
type Model struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	PublicID    string `gorm:"type:uuid;uniqueIndex"`
	Title       string
	Description string
	Status      string
	Priority    int
	DueDate     *time.Time
	Tags        string // Store as JSON or comma-separated
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Model) TableName() string {
	return "tasks"
}

// Convert domain entity → database model
func toModel(t *task.Entity) *Model {
	return &Model{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    t.Priority,
		DueDate:     t.DueDate,
		Tags:        strings.Join(t.Tags, ","),
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// Convert database model → domain entity
func toEntity(m *Model) *task.Entity {
	var tags []string
	if m.Tags != "" {
		tags = strings.Split(m.Tags, ",")
	}

	return &task.Entity{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      task.Status(m.Status),
		Priority:    m.Priority,
		DueDate:     m.DueDate,
		Tags:        tags,
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
