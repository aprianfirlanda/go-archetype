package taskhandler

import (
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	log         *logrus.Entry
	taskService portin.TaskService
}

func NewHandler(handlerLog *logrus.Entry, taskService portin.TaskService) *Handler {
	handlerLog = logging.WithComponent(handlerLog, "http.TaskHandler")
	return &Handler{
		log:         handlerLog,
		taskService: taskService,
	}
}
