package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
	"go-archetype/internal/application/task/command"

	"github.com/gofiber/fiber/v2"
)

// BulkDelete godoc
// @Summary      Bulk delete tasks
// @Tags         tasks
// @Accept       json
// @Security     JWTAuth
// @Param        request body taskreq.BulkDelete true "Bulk DeletePublicID Tasks"
// @Success      200  {object} response.Success{data=taskresp.BulkDelete}
// @Failure      400 {object} response.Error{errors=taskresp.BulkDeleteValidateError}
// @Router       /v1/api/tasks [delete]
func (h *Handler) BulkDelete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var req taskreq.BulkDelete
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
	res, err := h.taskService.BulkDelete(c.Context(), cmd)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.OK(res, rid))
}
