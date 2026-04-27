package db

import (
	"fmt"

	"go-archetype/internal/infrastructure/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenPostgres(cfg config.Database, logger *logrus.Entry) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode, cfg.TimeZone,
	)

	gormLogger := NewGormLogrusLogger(
		logger,
		cfg.LogLevel,
		WithSlowThreshold(cfg.SlowThreshold),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	logger.
		WithField("component", "infrastructure.db.postgres").
		Info("successfully connected to postgres")

	return db, nil
}
