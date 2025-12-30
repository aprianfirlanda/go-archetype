package task

import (
	"context"
	"errors"
	"go-archetype/internal/domain/identity"
	taskDomain "go-archetype/internal/domain/task"
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"
	"time"
)

type Service struct {
	taskRepository taskDomain.Repository
	uow            output.UnitOfWork
}

func NewService(uow output.UnitOfWork, taskRepository taskDomain.Repository) input.TaskService {
	return &Service{
		uow:            uow,
		taskRepository: taskRepository,
	}
}

func (s *Service) CreateTask(ctx context.Context, task *taskDomain.Entity) (string, error) {
	// Business validation
	if err := task.Validate(); err != nil {
		return "", err
	}

	// Set defaults
	task.PublicID = identity.NewPublicID()
	task.Status = taskDomain.StatusTodo
	task.Completed = false
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	return task.PublicID, s.taskRepository.Create(ctx, task)
}

func (s *Service) GetTaskByPublicID(ctx context.Context, publicID string) (*taskDomain.Entity, error) {
	if publicID == "" {
		return nil, errors.New("task publicID is required")
	}
	return s.taskRepository.FindByPublicID(ctx, publicID)
}

func (s *Service) ListTasks(ctx context.Context, filter taskDomain.ListFilter) ([]*taskDomain.Entity, int64, error) {
	return s.taskRepository.FindAll(ctx, filter)
}

func (s *Service) UpdateTask(ctx context.Context, task *taskDomain.Entity) error {
	// Validate
	if err := task.Validate(); err != nil {
		return err
	}

	// Check if task exists
	existing, err := s.taskRepository.FindByPublicID(ctx, task.PublicID)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("task not found")
	}

	task.UpdatedAt = time.Now()
	return s.taskRepository.Update(ctx, task)
}

func (s *Service) UpdateTaskStatus(ctx context.Context, publicID string, status taskDomain.Status) error {
	task, err := s.taskRepository.FindByPublicID(ctx, publicID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	if err := task.UpdateStatus(status); err != nil {
		return err
	}

	return s.taskRepository.Update(ctx, task)
}

func (s *Service) DeleteTaskByPublicID(ctx context.Context, publicID string) error {
	if publicID == "" {
		return errors.New("task publicID is required")
	}

	// Check if exists
	task, err := s.taskRepository.FindByPublicID(ctx, publicID)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}

	return s.taskRepository.DeleteByPublicID(ctx, publicID)
}

func (s *Service) BulkUpdateStatus(ctx context.Context, publicIDs []string, status taskDomain.Status) error {
	if len(publicIDs) == 0 {
		return errors.New("no task publicIDs provided")
	}

	if !status.IsValid() {
		return errors.New("invalid status")
	}

	// Start transaction
	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update each task with business validation
	for _, publicID := range publicIDs {
		if err := s.UpdateTaskStatus(ctx, publicID, status); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Service) BulkDelete(ctx context.Context, publicIDs []string) error {
	if len(publicIDs) == 0 {
		return errors.New("no task IDs provided")
	}

	tx, err := s.uow.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, publicID := range publicIDs {
		if err := s.DeleteTaskByPublicID(ctx, publicID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
