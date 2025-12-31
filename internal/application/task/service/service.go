package tasksvc

import (
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"
)

type Service struct {
	taskRepository portout.TaskRepository
	uow            portout.UnitOfWork
}

func NewService(uow portout.UnitOfWork, taskRepository portout.TaskRepository) portin.TaskService {
	return &Service{
		uow:            uow,
		taskRepository: taskRepository,
	}
}
