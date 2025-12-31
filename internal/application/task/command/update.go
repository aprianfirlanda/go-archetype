package taskcmd

import (
	"time"
)

type Update struct {
	PublicID    string
	Title       string
	Description string
	Priority    int
	DueDate     *time.Time
	Tags        []string
}
