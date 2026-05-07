package taskgorm

import (
	"context"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/output"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func componentLog(ctx context.Context) *logrus.Entry {
	return logging.ComponentLogger(logging.FromContext(ctx), "persistence.gorm.task.repository")
}

func New(db *gorm.DB) portout.TaskRepository {
	return &repository{db: db}
}
