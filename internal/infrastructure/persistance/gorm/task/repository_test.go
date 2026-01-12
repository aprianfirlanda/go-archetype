package taskgorm

import (
	"context"
	taskquery "go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/identity"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/infrastructure/testutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *testutil.TestDB

// 🔹 Runs ONCE for this package
func TestMain(m *testing.M) {
	ctx := context.Background()

	db, err := testutil.StartPostgres(ctx, "migrations")
	if err != nil {
		panic(err)
	}
	testDB = db

	code := m.Run()

	sqlDB, _ := testDB.DB.DB()
	_ = sqlDB.Close()

	os.Exit(code)
}

func TestTaskRepository_CRUD(t *testing.T) {
	ctx := context.Background()

	// Clean data for this test
	defer func() {
		_ = testutil.Truncate(testDB.DB, "tasks")
	}()

	repo := New(testDB.DB)

	now := time.Now()
	entity := &task.Entity{
		PublicID:    identity.NewPublicID(),
		Title:       "Test Task",
		Description: "Testing repository",
		Status:      task.StatusTodo,
		Priority:    1,
		Tags:        []string{"test", "repo"},
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// CREATE
	err := repo.Create(ctx, entity)
	require.NoError(t, err)

	// FIND
	found, err := repo.FindByPublicID(ctx, entity.PublicID)
	require.NoError(t, err)
	assert.Equal(t, entity.Title, found.Title)

	// UPDATE
	entity.Title = "Updated Title"
	err = repo.UpdateByPublicID(ctx, entity)
	require.NoError(t, err)

	updated, err := repo.FindByPublicID(ctx, entity.PublicID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Title", updated.Title)

	// DELETE
	err = repo.DeleteByPublicID(ctx, entity.PublicID)
	require.NoError(t, err)

	_, err = repo.FindByPublicID(ctx, entity.PublicID)
	assert.Equal(t, task.ErrNotFound, err)
}

func TestTaskRepository_FindAll(t *testing.T) {
	ctx := context.Background()

	defer func() {
		_ = testutil.Truncate(testDB.DB, "tasks")
	}()

	repo := New(testDB.DB)

	// Seed data
	for i := 1; i <= 5; i++ {
		_ = repo.Create(ctx, &task.Entity{
			PublicID:  identity.NewPublicID(),
			Title:     "Task",
			Status:    task.StatusTodo,
			Priority:  i,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	filter := taskquery.ListFilter{
		Page:  1,
		Limit: 2,
	}

	list, total, err := repo.FindAll(ctx, filter)
	require.NoError(t, err)

	assert.Equal(t, int64(5), total)
	assert.Len(t, list, 2)
}
