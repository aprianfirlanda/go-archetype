package gorm

import (
	"context"
	"go-archetype/internal/infrastructure/database"
	"go-archetype/internal/ports/outbound"

	"gorm.io/gorm"
)

type Pinger struct {
	DB *gorm.DB
}

func NewPinger(db *gorm.DB) outbound.DBPinger {
	return &Pinger{DB: db}
}

func (p Pinger) Ping(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return database.Ping(ctx, sqlDB)
}
