package task

import "time"

type Task struct {
	ID          string
	Title       string
	Description string
	Status      Status
	Priority    int
	DueDate     *time.Time
	Tags        []string
	Completed   bool
	CreatedAt   time.Time
}
