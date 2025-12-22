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

// Create godoc
// @Summary      Create a task
// @Description  Create a new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        request body request.CreateTask true "Create Task Request"
// @Success      201 {object} response.Simple{data=response.Task}
// @Failure      400 {object} response.ErrorResponse{errors=response.CreateTaskValidateError}
// @Router       /api/tasks [post]
func (h *TaskHandler) Create(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.CreateTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
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
	return c.Status(fiber.StatusCreated).JSON(response.Simple{Data: resp})
}

// GetByID godoc
// @Summary      Get task by ID
// @Description  Retrieve a single task by its ID
// @Tags         tasks
// @Produce      json
// @Security     JWTAuth
// @Security     ApiKeyAuth
// @Param        id   path     string  true  "Task ID"
// @Success      200 {object} response.Simple{data=response.Task}
// @Failure      400  {object} response.ErrorResponse
// @Failure      404  {object} response.ErrorResponse
// @Router       /api/tasks/{id} [get]
func (h *TaskHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "task ID is required")
	}

	return c.JSON(response.Simple{Data: response.Task{
		ID:     id,
		Title:  "Demo Task",
		Status: task.StatusInProgress,
	}})
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
// @Param        status   query string  false "Task status"
// @Param        priority query int     false "Task priority"
// @Success      200 {object} response.Paginate{data=[]response.Task}
// @Failure      400 {object} response.ErrorResponse{errors=response.ListTasksValidateError}
// @Router       /api/tasks [get]
func (h *TaskHandler) List(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var q request.ListTasks
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(q)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
	}

	q.Normalize()

	log.WithFields(logrus.Fields{
		"page":   q.Page,
		"limit":  q.Limit,
		"search": q.Search,
		"status": q.Status,
	}).Info("list tasks")

	return c.JSON(response.Paginate{
		Data: []response.Task{},
		Meta: response.Meta{
			Page:       q.Page,
			PerPage:    q.Limit,
			TotalItems: 0,
			TotalPages: 0,
			HasNext:    false,
			HasPrev:    false,
		},
	})
}

// Update godoc
// @Summary      Update task
// @Description  Replace a task completely
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        id      path string true "Task ID"
// @Param        request body request.UpdateTask true "Update Task Request"
// @Success      204
// @Failure      400 {object} response.ErrorResponse{errors=response.UpdateTaskValidateError}
// @Router       /api/tasks/{id} [put]
func (h *TaskHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "task ID is required")
	}

	var req request.UpdateTask
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateStatus godoc
// @Summary      Update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        id      path string true "Task ID"
// @Param        request body request.UpdateTaskStatus true "Update Status Request"
// @Success      204
// @Failure      400 {object} response.ErrorResponse{errors=response.UpdateTaskStatusValidateError}
// @Router       /api/tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "task ID is required")
	}

	var req request.UpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
	}

	log.WithFields(logrus.Fields{
		"task_id": id,
		"status":  req.Status,
	}).Info("update task status")

	// later:
	// err := h.updateTaskStatus.Execute(ctx, id, req.Status)

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
// @Failure      400 {object} response.ErrorResponse{errors=response.BulkUpdateTaskStatusValidateError}
// @Router       /api/tasks/status [patch]
func (h *TaskHandler) BulkUpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.BulkUpdateTaskStatus
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
	}

	log.WithFields(logrus.Fields{
		"task_ids": req.IDs,
		"status":   req.Status,
	}).Info("bulk update task status")

	// later:
	// err := h.updateStatuses.Execute(c.Context(), req.IDs, req.Status)

	return c.SendStatus(fiber.StatusNoContent)
}

// Delete godoc
// @Summary      Delete task
// @Tags         tasks
// @Security     JWTAuth
// @Param        id path string true "Task ID"
// @Success      204
// @Failure      400 {object} response.ErrorResponse
// @Router       /api/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	id := c.Params("id")
	if id == "" {
		return fiber.NewError(fiber.StatusBadRequest, "task ID is required")
	}

	log.WithField("task_id", id).Info("delete task")

	// later:
	// err := h.deleteTask.Execute(ctx, id)

	return c.SendStatus(fiber.StatusNoContent)
}

// BulkDelete godoc
// @Summary      Bulk delete tasks
// @Tags         tasks
// @Accept       json
// @Security     JWTAuth
// @Param        request body request.BulkDeleteTasks true "Bulk Delete Tasks"
// @Success      204
// @Failure      400 {object} response.ErrorResponse{errors=response.BulkDeleteTasksValidateError}
// @Router       /api/tasks [delete]
func (h *TaskHandler) BulkDelete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	var req request.BulkDeleteTasks
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	errors, err := validation.ValidateStruct(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			response.ErrorResponse{
				Message:   "validation failed",
				Errors:    errors,
				RequestID: httpctx.GetRequestID(c),
			},
		)
	}

	log.WithField("task_ids", req.IDs).Info("bulk delete tasks")

	// later (usecase):
	// err := h.deleteTasks.Execute(c.Context(), req.IDs)
	// if err != nil { return err }

	return c.SendStatus(fiber.StatusNoContent)
}
