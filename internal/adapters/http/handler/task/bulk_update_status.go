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

// BulkUpdateStatus godoc
// @Summary      Bulk update task status
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        request body taskreq.BulkUpdateStatus true "Bulk Update Status"
// @Success      200  {object} response.Success{data=taskresp.BulkUpdateStatus}
// @Failure      400 {object} response.Error{errors=taskresp.BulkUpdateStatusValidateError}
// @Router       /v1/api/tasks/status [patch]
func (h *Handler) BulkUpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var req taskreq.BulkUpdateStatus
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

	cmd := taskcmd.BulkUpdateStatus{
		PublicIDs: req.IDs,
		Status:    task.Status(req.Status),
	}
	res, err := h.taskService.BulkUpdateStatus(c.Context(), cmd)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.OK(res, rid))
}
