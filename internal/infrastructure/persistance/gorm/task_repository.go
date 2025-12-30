package gorm

import (
	"context"
	"errors"
	"go-archetype/internal/domain/task"
	"strings"
	"time"

	"gorm.io/gorm"
)

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) task.Repository {
	return &taskRepository{db: db}
}

// TaskModel is the GORM model (infrastructure concern)
type TaskModel struct {
	ID          int64  `gorm:"primaryKey;autoIncrement"`
	PublicID    string `gorm:"type:uuid;uniqueIndex"`
	Title       string
	Description string
	Status      string
	Priority    int
	DueDate     *time.Time
	Tags        string // Store as JSON or comma-separated
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TaskModel) TableName() string {
	return "tasks"
}

// Convert domain entity → database model
func toModel(t *task.Entity) *TaskModel {
	return &TaskModel{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Status:      string(t.Status),
		Priority:    t.Priority,
		DueDate:     t.DueDate,
		Tags:        strings.Join(t.Tags, ","),
		Completed:   t.Completed,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

// Convert database model → domain entity
func toEntity(m *TaskModel) *task.Entity {
	var tags []string
	if m.Tags != "" {
		tags = strings.Split(m.Tags, ",")
	}

	return &task.Entity{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Status:      task.Status(m.Status),
		Priority:    m.Priority,
		DueDate:     m.DueDate,
		Tags:        tags,
		Completed:   m.Completed,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *taskRepository) Create(ctx context.Context, t *task.Entity) error {
	model := toModel(t)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *taskRepository) FindByPublicID(ctx context.Context, publicID string) (*task.Entity, error) {
	var model TaskModel
	err := r.db.WithContext(ctx).First(&model, "public_id = ?", publicID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, task.ErrNotFound
		}
		return nil, err
	}
	return toEntity(&model), nil
}

func (r *taskRepository) FindAll(ctx context.Context, filter task.ListFilter) ([]*task.Entity, int64, error) {
	query := r.db.WithContext(ctx).Model(&TaskModel{})

	// Apply filters
	if filter.Search != "" {
		query = query.Where("title LIKE ? OR description LIKE ?",
			"%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", string(filter.Status))
	}
	if filter.Priority != nil {
		query = query.Where("priority = ?", *filter.Priority)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := (filter.Page - 1) * filter.Limit
	query = query.Offset(offset).Limit(filter.Limit)

	// Fetch
	var models []TaskModel
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to entities
	tasks := make([]*task.Entity, len(models))
	for i, m := range models {
		tasks[i] = toEntity(&m)
	}

	return tasks, total, nil
}

func (r *taskRepository) Update(ctx context.Context, t *task.Entity) error {
	model := toModel(t)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *taskRepository) DeleteByPublicID(ctx context.Context, publicID string) error {
	return r.db.WithContext(ctx).Delete(&TaskModel{}, "public_id = ?", publicID).Error
}

func (r *taskRepository) BulkUpdateStatus(ctx context.Context, publicIDs []string, status task.Status) error {
	return r.db.WithContext(ctx).
		Model(&TaskModel{}).
		Where("public_id IN ?", publicIDs).
		Updates(map[string]interface{}{
			"status":     string(status),
			"completed":  status == task.StatusDone,
			"updated_at": time.Now(),
		}).Error
}

func (r *taskRepository) BulkDelete(ctx context.Context, publicIDs []string) error {
	return r.db.WithContext(ctx).Delete(&TaskModel{}, "public_id IN ?", publicIDs).Error
}
