package db

import (
	"fmt"

	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/infrastructure/logging"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// OpenPostgres opens a PostgreSQL connection using GORM,
func OpenPostgres(cfg config.Database, logger *logrus.Entry) (*gorm.DB, error) {
	log := logging.WithComponent(logger, "infrastructure.db.postgres")

	gormCfg := &gorm.Config{}
	if lvl := parseLogLevel(cfg.LogLevel); lvl != nil {
		gormCfg.Logger = gormlogger.Default.LogMode(*lvl)
	}

	db, err := gorm.Open(postgres.Open(buildDSN(cfg)), gormCfg)
	if err != nil {
		log.WithError(err).Error("failed to open postgres")
		return nil, err
	}

	log.Info("successfully connected to postgres")

	return db, nil
}

// buildPostgresDSN constructs a standard PostgreSQL DSN for GORM.
func buildDSN(cfg config.Database) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.Port,
		cfg.SSLMode,
		cfg.TimeZone,
	)
}

// parseLogLevel converts a string into a gormlogger.LogLevel pointer.
// Returns nil if the level is empty/unknown (caller should use default(warn)).
func parseLogLevel(level string) *gormlogger.LogLevel {
	if level == "" {
		return nil
	}

	switch level {
	case "silent":
		l := gormlogger.Silent
		return &l
	case "error":
		l := gormlogger.Error
		return &l
	case "warn":
		l := gormlogger.Warn
		return &l
	case "info":
		l := gormlogger.Info
		return &l
	default:
		return nil
	}
}
