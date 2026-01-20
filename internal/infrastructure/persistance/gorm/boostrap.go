package gorminfra

import (
	"context"
	"fmt"

	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/infrastructure/db"
	"go-archetype/internal/infrastructure/logging"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitPostgres(
	cfg config.Database,
	logger *logrus.Entry,
	autoMigrateModels []any,
) (*gorm.DB, error) {
	log := logging.WithComponent(logger, "infrastructure.persistence.gorm.boostrap")

	// 1. Open connection
	gormDB, err := db.OpenPostgres(cfg, logger)
	if err != nil {
		return nil, err
	}

	// 2. Configure pool
	sqlDB, err := gormDB.DB()
	if err != nil {
		log.WithError(err).Error("failed to get sql.DB")
		return nil, err
	}
	db.ConfigurePool(sqlDB, cfg)

	// 3. Ping
	if err := db.Ping(context.Background(), sqlDB); err != nil {
		log.WithError(err).Error("failed to ping postgres")
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	// 4. Auto-migrate (GORM concern)
	if len(autoMigrateModels) > 0 {
		if err := gormDB.AutoMigrate(autoMigrateModels...); err != nil {
			log.WithError(err).Error("auto-migrate failed")
			return nil, err
		}
	}

	return gormDB, nil
}
