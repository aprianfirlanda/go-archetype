package gorm

import (
	"context"
	"go-archetype/internal/infrastructure/database"

	"gorm.io/gorm"
)

type Pinger struct {
	DB *gorm.DB
}

func (p Pinger) Ping(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return database.Ping(ctx, sqlDB)
}
