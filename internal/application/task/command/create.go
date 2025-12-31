package taskcmd

import (
	"time"
)

type Create struct {
	Title       string
	Description string
	Priority    int
	DueDate     *time.Time
	Tags        []string
}
