package task

import "context"

// Repository defines what task domain needs from storage
type Repository interface {
	Create(ctx context.Context, task *Entity) error
	FindByPublicID(ctx context.Context, publicID string) (*Entity, error)
	FindAll(ctx context.Context, filter ListFilter) ([]*Entity, int64, error)
	Update(ctx context.Context, task *Entity) error
	DeleteByPublicID(ctx context.Context, publicID string) error
	BulkUpdateStatus(ctx context.Context, publicIDs []string, status Status) error
	BulkDelete(ctx context.Context, publicIDs []string) error
}

type ListFilter struct {
	Search   string
	Status   Status
	Priority *int
	Page     int
	Limit    int
}
