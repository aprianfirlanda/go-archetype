package taskhandlermq

import "go-archetype/internal/ports/input"

type Handler struct {
	service portin.TaskService
}

func New(service portin.TaskService) *Handler {
	return &Handler{
		service: service,
	}
}
