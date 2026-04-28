package taskhandlermq

import portin "go-archetype/internal/ports/input"

type Handler struct {
	service portin.TaskService
}

func New(service portin.TaskService) *Handler {
	return &Handler{
		service: service,
	}
}
