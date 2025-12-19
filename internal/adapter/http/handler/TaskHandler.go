package handler

import (
	httpctx "go-archetype/internal/adapter/http/context"
	"go-archetype/internal/adapter/http/dto/request"
	"go-archetype/internal/adapter/http/dto/response"
	"go-archetype/internal/adapter/http/validation"
	"go-archetype/internal/domain/task"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	log *logrus.Entry
}

func NewTaskHandler(log *logrus.Entry) *TaskHandler {
	return &TaskHandler{
		log: log,
	}
}

func (h *TaskHandler) Create(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.CreateTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := validation.ValidateStruct(req); err != nil {
		log.WithError(err).Warn("validation failed")
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp := response.Task{
		ID:          "task-123",
		Title:       req.Title,
		Description: req.Description,
		Status:      task.StatusTodo,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	log.WithField("task_id", resp.ID).Info("task created")
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	return c.JSON(response.Task{
		ID:     id,
		Title:  "Demo Task",
		Status: task.StatusInProgress,
	})
}

func (h *TaskHandler) List(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var q request.ListTasks
	if err := c.QueryParser(&q); err != nil {
		log.WithError(err).Warn("invalid query params")
		return fiber.ErrBadRequest
	}

	if err := validation.ValidateStruct(q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	q.Normalize()

	log.WithFields(logrus.Fields{
		"page":   q.Page,
		"limit":  q.Limit,
		"search": q.Search,
		"status": q.Status,
	}).Info("list tasks")

	return c.JSON(fiber.Map{
		"data":  []response.Task{},
		"page":  q.Page,
		"limit": q.Limit,
		"total": 0,
	})
}

func (h *TaskHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req request.UpdateTask

	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	if err := validation.ValidateStruct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"id":     id,
		"status": "updated",
	})
}

func (h *TaskHandler) UpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	id := c.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	var req request.UpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Warn("invalid request body")
		return fiber.ErrBadRequest
	}

	log.WithFields(logrus.Fields{
		"task_id": id,
		"status":  req.Status,
	}).Info("update task status")

	// later:
	// err := h.updateTaskStatus.Execute(ctx, id, req.Status)

	return c.JSON(fiber.Map{
		"id":     id,
		"status": req.Status,
	})
}

func (h *TaskHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.BulkUpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Warn("invalid bulk status update body")
		return fiber.ErrBadRequest
	}

	log.WithFields(logrus.Fields{
		"task_ids": req.IDs,
		"status":   req.Status,
	}).Info("bulk update task status")

	// later:
	// err := h.updateStatuses.Execute(c.Context(), req.IDs, req.Status)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	id := c.Params("id")
	if id == "" {
		return fiber.ErrBadRequest
	}

	log.WithField("task_id", id).Info("delete task")

	// later:
	// err := h.deleteTask.Execute(ctx, id)

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TaskHandler) BulkDelete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.BulkDeleteTasks
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Warn("invalid bulk delete request body")
		return fiber.ErrBadRequest
	}

	if len(req.IDs) == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "ids must not be empty")
	}

	log.WithField("task_ids", req.IDs).Info("bulk delete tasks")

	// later (usecase):
	// err := h.deleteTasks.Execute(c.Context(), req.IDs)
	// if err != nil { return err }

	return c.SendStatus(fiber.StatusNoContent)
}
