package db

import (
	"context"
	"database/sql"
	"fmt"
	"go-archetype/internal/config"
	"go-archetype/internal/logging"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// InitPostgres opens a PostgreSQL connection using GORM,
// configures connection pool & logging, optionally runs AutoMigrate.
func InitPostgres(dbCfg config.Database, logger *logrus.Entry, autoMigrateModels []any) (*gorm.DB, error) {
	log := logging.WithComponent(logger, "db.gorm")

	dsn := buildPostgresDSN(dbCfg)

	gormCfg := &gorm.Config{}
	if lvl := parseLogLevel(dbCfg.LogLevel); lvl != nil {
		gormCfg.Logger = gormlogger.Default.LogMode(*lvl)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		log.WithError(err).Error("failed to open postgres with gorm")
		return nil, fmt.Errorf("open postgres with gorm: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.WithError(err).Error("failed to get underlying *sql.DB")
		return nil, fmt.Errorf("get underlying *sql.DB: %w", err)
	}

	configureConnectionPool(sqlDB, dbCfg)

	if err := Ping(context.Background(), sqlDB); err != nil {
		log.WithError(err).Error("failed to ping postgres")
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	// Optionally, run AutoMigrate (code-first schema management).
	if len(autoMigrateModels) > 0 {
		if err := db.AutoMigrate(autoMigrateModels...); err != nil {
			log.WithError(err).Error("failed to auto-migrate")
			return nil, fmt.Errorf("auto-migrate: %w", err)
		}
	}

	log.Info("successfully connected to postgres")

	return db, nil
}

// buildPostgresDSN constructs a standard PostgreSQL DSN for GORM.
func buildPostgresDSN(dbCfg config.Database) string {
	// Example:
	// host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable TimeZone=UTC
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbCfg.Host,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.Name,
		dbCfg.Port,
		dbCfg.SSLMode,
		dbCfg.TimeZone,
	)
}

// parseLogLevel converts a string into a gormlogger.LogLevel pointer.
// Returns nil if the level is empty/unknown (caller should use default(warn)).
func parseLogLevel(level string) *gormlogger.LogLevel {
	if level == "" {
		return nil
	}

	var lvl gormlogger.LogLevel
	switch level {
	case "silent", "SilENT", "SILENT":
		lvl = gormlogger.Silent
	case "error", "Error", "ERROR":
		lvl = gormlogger.Error
	case "warn", "Warn", "WARN", "warning", "Warning", "WARNING":
		lvl = gormlogger.Warn
	case "info", "Info", "INFO":
		lvl = gormlogger.Info
	default:
		// Unknown string: fall back to nil -> use GORM default(warn).
		return nil
	}
	return &lvl
}

func configureConnectionPool(sqlDB *sql.DB, dbCfg config.Database) {
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbCfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(dbCfg.ConnMaxIdleTime)
}

// Ping checks database reachability with a timeout.
// Exported so it can be reused by other packages (e.g., readiness checks).
func Ping(ctx context.Context, sqlDB *sql.DB) error {
	if sqlDB == nil {
		return fmt.Errorf("sqlDB is nil")
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}
	return nil
}
