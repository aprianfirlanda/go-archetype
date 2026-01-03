package tasksvc

import (
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"
)

type service struct {
	taskRepository portout.TaskRepository
	uow            portout.UnitOfWork
}

func New(uow portout.UnitOfWork, taskRepository portout.TaskRepository) portin.TaskService {
	return &service{
		uow:            uow,
		taskRepository: taskRepository,
	}
}
