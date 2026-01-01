package taskhandler

import (
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"

	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	log         *logrus.Entry
	taskService portin.TaskService
}

func NewTaskHandler(handlerLog *logrus.Entry, taskService portin.TaskService) *TaskHandler {
	handlerLog = logging.WithComponent(handlerLog, "http.TaskHandler")
	return &TaskHandler{
		log:         handlerLog,
		taskService: taskService,
	}
}
