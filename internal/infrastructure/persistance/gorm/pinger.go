package gorm

import (
	"context"
	"go-archetype/internal/infrastructure/db"
	"go-archetype/internal/ports"

	"gorm.io/gorm"
)

type Pinger struct {
	db *gorm.DB
}

func NewPinger(db *gorm.DB) ports.DBPinger {
	return &Pinger{db: db}
}

func (p *Pinger) Ping(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Ping(ctx, sqlDB)
}
