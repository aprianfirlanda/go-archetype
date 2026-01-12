package testutil

import (
	"context"
	"fmt"
	"go-archetype/internal/infrastructure/config"
	infragorm "go-archetype/internal/infrastructure/persistance/gorm"
	"go-archetype/internal/infrastructure/persistance/gorm/migrate"
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

	db, err := infragorm.InitPostgres(cfg, logger, nil)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
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

	return &TestDB{DB: db}, nil
}
