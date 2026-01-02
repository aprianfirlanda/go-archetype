package taskgorm

import (
	"go-archetype/internal/domain/task"
	"strings"
	"time"
)

// Task is the GORM model (infrastructure concern)
type Task struct {
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

func (Task) TableName() string {
	return "tasks"
}

// Convert domain entity → database model
func toModel(t *task.Entity) *Task {
	return &Task{
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
func toEntity(m *Task) *task.Entity {
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
