package tasksvc

import (
	"context"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"
	"go-archetype/internal/ports/output"

	"github.com/sirupsen/logrus"
)

type service struct {
	taskRepository portout.TaskRepository
	uow            portout.UnitOfWork
	publisher      portout.MessagePublisher
}

func componentLog(ctx context.Context) *logrus.Entry {
	return logging.ComponentLogger(logging.FromContext(ctx), "application.task.service")
}

func New(
	uow portout.UnitOfWork,
	taskRepository portout.TaskRepository,
	publisher portout.MessagePublisher,
) portin.TaskService {
	return &service{
		uow:            uow,
		taskRepository: taskRepository,
		publisher:      publisher,
	}
}
