package migrate

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

type GooseMigrator struct {
	db  *sql.DB
	dir string
}

func NewGooseMigrator(db *sql.DB, migrationsDir string) *GooseMigrator {
	return &GooseMigrator{
		db:  db,
		dir: migrationsDir,
	}
}

// Up applies all pending migrations
func (g *GooseMigrator) Up(_ context.Context) error {
	return goose.Up(g.db, g.dir)
}

// UpTo migrates up to a specific version
func (g *GooseMigrator) UpTo(_ context.Context, version int64) error {
	return goose.UpTo(g.db, g.dir, version)
}

// Down rolls back the last migration
func (g *GooseMigrator) Down(_ context.Context) error {
	return goose.Down(g.db, g.dir)
}

// DownTo rolls back migrations down to a specific version
func (g *GooseMigrator) DownTo(_ context.Context, version int64) error {
	return goose.DownTo(g.db, g.dir, version)
}

// Status prints migration status
func (g *GooseMigrator) Status(_ context.Context) error {
	return goose.Status(g.db, g.dir)
}

// Version returns current migration version
func (g *GooseMigrator) Version(_ context.Context) (int64, error) {
	return goose.GetDBVersion(g.db)
}

// Create creates a new migration file
func (g *GooseMigrator) Create(name string, migrationType string) error {
	return goose.Create(nil, g.dir, name, migrationType)
}
