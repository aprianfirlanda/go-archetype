package taskhandler

import (
	"errors"
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"

	"github.com/gofiber/fiber/v2"
)

// Update godoc
// @Summary      Update task
// @Description  Replace a task completely
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        public_id   path     string  true  "Entity Public ID"
// @Param        request body taskreq.Update true "Update Entity Request"
// @Success      204
// @Failure      400  {object} response.Error{errors=taskresp.UpdateValidateError}
// @Failure      404  {object} response.Error
// @Router       /v1/api/tasks/{public_id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task publicID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	var req taskreq.Update
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
