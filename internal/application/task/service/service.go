package tasksvc

import (
	taskport "go-archetype/internal/application/task/port"
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"
)

type Service struct {
	taskRepository taskport.Repository
	uow            portout.UnitOfWork
}

func NewService(uow portout.UnitOfWork, taskRepository taskport.Repository) portin.TaskService {
	return &Service{
		uow:            uow,
		taskRepository: taskRepository,
	}
}
