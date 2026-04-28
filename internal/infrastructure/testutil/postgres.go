package testutil

import (
	"context"
	"fmt"
	"go-archetype/internal/adapters/persistence/gorm/migrate"
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/infrastructure/db"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TestDB struct {
	DB *gorm.DB
}

func StartPostgres(ctx context.Context, migrationsDir string) (*TestDB, error) {
	logger := logrus.New().WithField("component", "test")

	cfg := config.Database{
		Host:     "localhost",
		Port:     5432,
		User:     "app",
		Password: "change_me",
		Name:     "app",
		SSLMode:  "disable",
		TimeZone: "UTC",
	}

	dbConn, err := db.NewPostgres(cfg, logger, nil)
	if err != nil {
		return nil, err
	}

	sqlDB, err := dbConn.DB()
	if err != nil {
		return nil, err
	}

	root, err := ProjectRoot()
	if err != nil {
		return nil, err
	}

	migrationsDir = filepath.Join(root, migrationsDir)
	goose := migrate.NewGooseMigrator(sqlDB, migrationsDir)
	if err := goose.Up(ctx); err != nil {
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return &TestDB{DB: dbConn}, nil
}
