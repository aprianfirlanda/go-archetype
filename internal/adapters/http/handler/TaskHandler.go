package handler

import (
	"errors"
	httpctx "go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/converter"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"
	"go-archetype/internal/infrastructure/logging"
	"go-archetype/internal/ports/input"

	"github.com/gofiber/fiber/v2"
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

	cmd := taskcmd.Create{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	}
	publicID, err := h.taskService.Create(c.Context(), cmd)
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
// @Param        public_id   path     string  true  "Entity Public ID"
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

	taskEntity, err := h.taskService.GetByPublicID(c.Context(), publicID)
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

	tasks, total, err := h.taskService.List(c.Context(), filter)
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
// @Param        public_id   path     string  true  "Entity Public ID"
// @Param        request body request.UpdateTask true "Update Entity Request"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.UpdateTaskValidateError}
// @Failure      404  {object} response.Error
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

	cmd := taskcmd.Update{
		PublicID:    publicID,
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	}
	err = h.taskService.Update(c.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(response.FailMessage("task not found", rid))
		default:
			log.WithError(err).Error("failed to update task")
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailMessage("failed to update task", rid))
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateStatus godoc
// @Summary      Update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        public_id   path     string  true  "Entity Public ID"
// @Param        request body request.UpdateTaskStatus true "Update Status Request"
// @Success      204
// @Failure      400 {object} response.Error{errors=response.UpdateTaskStatusValidateError}
// @Failure      404 {object} response.Error
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

	status := task.Status(req.Status)
	if status.IsValid() {
		log.WithField("status", status).Info("invalid status")
		return c.Status(fiber.StatusBadRequest).JSON(response.OKMessage("invalid status", rid))
	}
	cmd := taskcmd.UpdateStatus{
		PublicID: publicID,
		Status:   status,
	}
	err = h.taskService.UpdateStatusSingle(c.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, task.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(response.FailMessage("task not found", rid))
		default:
			log.WithError(err).Error("failed to update task status")
			return c.Status(fiber.StatusInternalServerError).JSON(response.FailMessage("failed to update task status", rid))
		}
	}

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

	status := task.Status(req.Status)
	if status.IsValid() {
		log.WithField("status", status).Info("invalid status")
		return c.Status(fiber.StatusBadRequest).JSON(response.OKMessage("invalid status", rid))
	}
	cmd := taskcmd.BulkUpdateStatus{
		PublicIDs: req.IDs,
		Status:    status,
	}
	err = h.taskService.BulkUpdateStatus(c.Context(), cmd)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeletePublicID godoc
// @Summary      DeletePublicID task
// @Tags         tasks
// @Security     JWTAuth
// @Param        public_id   path     string  true  "Entity Public ID"
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

	if err := h.taskService.DeleteByPublicID(c.Context(), publicID); err != nil {
		return err
	}

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

	cmd := taskcmd.BulkDelete{
		PublicIDs: req.IDs,
	}
	if err := h.taskService.BulkDelete(c.Context(), cmd); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
