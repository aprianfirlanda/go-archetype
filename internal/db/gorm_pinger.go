package db

import (
	"context"

	"gorm.io/gorm"
)

type GormPinger struct {
	DB *gorm.DB
}

func (p GormPinger) Ping(ctx context.Context) error {
	sqlDB, err := p.DB.DB()
	if err != nil {
		return err
	}
	return Ping(ctx, sqlDB)
}
