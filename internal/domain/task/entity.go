package task

import (
	"errors"
	"go-archetype/internal/domain/identity"
	"time"
)

type Entity struct {
	ID          int64
	PublicID    string
	Title       string
	Description string
	Status      Status
	Priority    int
	DueDate     *time.Time
	Tags        []string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func New(
	title string,
	description string,
	priority int,
	dueDate *time.Time,
	tags []string,
) *Entity {
	now := time.Now()

	return &Entity{
		PublicID:    identity.NewPublicID(),
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     dueDate,
		Tags:        tags,
		Status:      StatusTodo,
		Completed:   false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Validate contains business rules
func (t *Entity) Validate() error {
	if t.Title == "" {
		return errors.New("title is required")
	}
	if len(t.Title) < 3 {
		return errors.New("title must be at least 3 characters")
	}
	if t.Priority < 1 || t.Priority > 5 {
		return errors.New("priority must be between 1 and 5")
	}
	if len(t.Tags) > 10 {
		return errors.New("maximum 10 tags allowed")
	}
	return nil
}

func (t *Entity) Update(
	title string,
	description string,
	priority int,
	dueDate *time.Time,
	tags []string,
) {
	t.Title = title
	t.Description = description
	t.Priority = priority
	t.DueDate = dueDate
	t.Tags = tags
	t.UpdatedAt = time.Now()
}

func (t *Entity) Complete() error {
	if t.Completed {
		return errors.New("task is already completed")
	}
	t.Completed = true
	t.Status = StatusDone
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Entity) UpdateStatus(status Status) error {
	if !status.IsValid() {
		return errors.New("invalid status")
	}

	if t.Status == StatusDone && status != StatusDone {
		return errors.New("cannot reopen completed task")
	}

	t.Status = status
	if status == StatusDone {
		t.Completed = true
	}
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Entity) IsOverdue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return t.DueDate.Before(time.Now())
}
