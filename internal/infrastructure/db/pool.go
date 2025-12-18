package db

import (
	"database/sql"

	"go-archetype/internal/infrastructure/config"
)

func ConfigurePool(sqlDB *sql.DB, cfg config.Database) {
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
}
