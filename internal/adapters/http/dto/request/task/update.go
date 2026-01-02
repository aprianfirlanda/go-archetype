package taskreq

import (
	"time"
)

type Update struct {
	Title       string     `json:"title" validate:"required,min=3"`
	Description string     `json:"description" validate:"max=500"`
	Priority    int        `json:"priority" validate:"required,min=1,max=5"`
	DueDate     *time.Time `json:"due_date"`
	Tags        []string   `json:"tags" validate:"max=10,dive,min=1"`
}
