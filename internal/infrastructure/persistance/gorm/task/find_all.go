package taskgorm

import (
	"context"
	"go-archetype/internal/application/task/query"
	"go-archetype/internal/domain/task"
)

func (r *repository) FindAll(ctx context.Context, filter taskquery.ListFilter) ([]*task.Entity, int64, error) {
	query := r.db.WithContext(ctx).Model(&Model{})

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
	var models []Model
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
