package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/response"

	"github.com/gofiber/fiber/v2"
)

// DeletePublicID godoc
// @Summary      DeletePublicID task
// @Tags         tasks
// @Security     JWTAuth
// @Param        public_id   path     string  true  "Entity Public ID"
// @Success      204
// @Failure      400 {object} response.Error
// @Router       /v1/api/tasks/{public_id} [delete]
func (h *Handler) DeletePublicID(c *fiber.Ctx) error {
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
