package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
	"go-archetype/internal/application/task/command"
	"go-archetype/internal/domain/task"

	"github.com/gofiber/fiber/v2"
)

// UpdateStatus godoc
// @Summary      Update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        public_id   path     string  true  "Entity Public ID"
// @Param        request body taskreq.UpdateStatus true "Update Status Request"
// @Success      204
// @Failure      400 {object} response.Error{errors=taskresp.UpdateStatusValidateError}
// @Failure      404 {object} response.Error
// @Router       /v1/api/tasks/{public_id}/status [patch]
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task publicID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	var req taskreq.UpdateStatus
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
	if !status.IsValid() {
		log.WithField("status", status).Info("invalid status")
		return c.Status(fiber.StatusBadRequest).JSON(response.OKMessage("invalid status", rid))
	}
	cmd := taskcmd.UpdateStatus{
		PublicID: publicID,
		Status:   status,
	}
	err = h.taskService.UpdateStatus(c.UserContext(), cmd)
	if err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
