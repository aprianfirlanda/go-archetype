package gorminfra

import (
	"context"
	"go-archetype/internal/infrastructure/db"
	"go-archetype/internal/ports/output"

	"gorm.io/gorm"
)

type Pinger struct {
	db *gorm.DB
}

func NewPinger(db *gorm.DB) portout.DBPinger {
	return &Pinger{db: db}
}

func (p *Pinger) Ping(ctx context.Context) error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return db.Ping(ctx, sqlDB)
}
