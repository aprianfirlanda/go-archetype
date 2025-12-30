package handler

import (
	"errors"
	httpctx "go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/converter"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type TaskHandler struct {
	log         *logrus.Entry
	taskService input.TaskService
}

func NewTaskHandler(handlerLog *logrus.Entry, taskService input.TaskService) *TaskHandler {
	handlerLog = logging.WithComponent(handlerLog, "http.TaskHandler")
	return &TaskHandler{
		log:         handlerLog,
		taskService: taskService,
	}
}

// Create godoc
// @Summary      Create a task
// @Description  Create a new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        request body request.CreateTask true "Create Entity Request"
// @Success      201 {object} response.Success{data=response.IDResponse}
// @Failure      400 {object} response.Error{errors=response.CreateTaskValidateError}
// @Failure      500 {object} response.Error
// @Router       /v1/api/tasks [post]
func (h *TaskHandler) Create(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var req request.CreateTask
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Error("failed to validate request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate request body", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	taskEntity := &task.Entity{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	}

	publicID, err := h.taskService.CreateTask(c.Context(), taskEntity)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Fail("failed to create task", err, rid))
	}

	return c.Status(fiber.StatusCreated).JSON(response.OK(response.IDResponse{ID: publicID}, rid))
}

// GetByPublicID godoc
// @Summary      Get task by ID
// @Description  Retrieve a single task by its ID
// @Tags         tasks
// @Produce      json
// @Security     JWTAuth
// @Security     ApiKeyAuth
// @Param        public_id   path     string  true  "Entity ID"
// @Success      200  {object} response.Success{data=task.Entity}
// @Failure      400  {object} response.Error
// @Failure      404  {object} response.Error
// @Router       /v1/api/tasks/{public_id} [get]
func (h *TaskHandler) GetByPublicID(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	taskEntity, err := h.taskService.GetTaskByPublicID(c.Context(), publicID)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(response.FailMessage("task not found", rid))
		default:
			log.WithError(err).Error("failed to get task")
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailMessage("failed to get task", rid))
		}
	}

	return c.JSON(response.OK(taskEntity, rid))
}

// List godoc
// @Summary      List tasks
// @Description  List tasks with pagination and filters
// @Tags         tasks
// @Produce      json
// @Security     JWTAuth
// @Param        page     query int     false "Page number"
// @Param        limit    query int     false "Page size"
// @Param        search   query string  false "Search keyword"
// @Param        status   query string  false "Entity status"
// @Param        priority query int     false "Entity priority"
// @Success      200 {object} response.Success{data=[]task.Entity, meta=response.PaginationMeta}
// @Failure      400 {object} response.Error{errors=response.ListTasksValidateError}
// @Router       /v1/api/tasks [get]
func (h *TaskHandler) List(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var q request.ListTasks
	if err := c.QueryParser(&q); err != nil {
		log.WithError(err).Error("failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse query params", rid))
	}

	fieldErrors, err := validation.ValidateStruct(q)
	if err != nil {
		log.WithError(err).Error("failed to validate query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate query params", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	q.Normalize()

	filter, err := converter.ToListFilter(q)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(response.FailMessage(err.Error(), rid))
	}

	tasks, total, err := h.taskService.ListTasks(c.Context(), filter)
	if err != nil {
		return err
	}

	meta := response.NewPaginationMeta(
		filter.Page,
		filter.Limit,
		total,
	)

	return c.Status(fiber.StatusOK).
		JSON(response.OKPaginate(tasks, meta, rid))
}

// Update godoc
// @Summary      Update task
// @Description  Replace a task completely
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        id      path string true "Entity ID"
// @Param        request body request.UpdateTask true "Update Entity Request"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.UpdateTaskValidateError}
// @Router       /v1/api/tasks/{public_id} [put]
func (h *TaskHandler) Update(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task publicID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	var req request.UpdateTask
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Error("failed to validate request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate request body", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateStatus godoc
// @Summary      Update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        id      path string true "Entity ID"
// @Param        request body request.UpdateTaskStatus true "Update Status Request"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.UpdateTaskStatusValidateError}
// @Router       /v1/api/tasks/{public_id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task publicID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	var req request.UpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Error("failed to validate request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate request body", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	log.WithFields(logrus.Fields{
		"task_id": publicID,
		"status":  req.Status,
	}).Info("update task status")

	// later:
	// err := h.updateTaskStatus.Execute(ctx, publicID, req.Status)

	return c.SendStatus(fiber.StatusNoContent)
}

// BulkUpdateStatus godoc
// @Summary      Bulk update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        request body request.BulkUpdateTaskStatus true "Bulk Update Status"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.BulkUpdateTaskStatusValidateError}
// @Router       /v1/api/tasks/status [patch]
func (h *TaskHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var req request.BulkUpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Error("failed to validate request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate request body", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	log.WithFields(logrus.Fields{
		"task_ids": req.IDs,
		"status":   req.Status,
	}).Info("bulk update task status")

	// later:
	// err := h.updateStatuses.Execute(c.Context(), req.IDs, req.Status)

	return c.SendStatus(fiber.StatusNoContent)
}

// DeletePublicID godoc
// @Summary      DeletePublicID task
// @Tags         tasks
// @Security     JWTAuth
// @Param        id path string true "Entity ID"
// @Success      204
// @Failure      400 {object} response.Error
// @Router       /v1/api/tasks/{public_id} [delete]
func (h *TaskHandler) DeletePublicID(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task publicID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	log.WithField("task_id", publicID).Info("delete task")

	// later:
	// err := h.deleteTask.Execute(ctx, publicID)

	return c.SendStatus(fiber.StatusNoContent)
}

// BulkDelete godoc
// @Summary      Bulk delete tasks
// @Tags         tasks
// @Accept       json
// @Security     JWTAuth
// @Param        request body request.BulkDeleteTasks true "Bulk DeletePublicID Tasks"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.BulkDeleteTasksValidateError}
// @Router       /v1/api/tasks [delete]
func (h *TaskHandler) BulkDelete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var req request.BulkDeleteTasks
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Error("failed to validate request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate request body", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	log.WithField("task_ids", req.IDs).Info("bulk delete tasks")

	// later (usecase):
	// err := h.deleteTasks.Execute(c.Context(), req.IDs)
	// if err != nil { return err }

	return c.SendStatus(fiber.StatusNoContent)
}
