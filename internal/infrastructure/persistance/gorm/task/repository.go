package taskgorm

import (
	"go-archetype/internal/ports/output"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) portout.TaskRepository {
	return &repository{db: db}
}
