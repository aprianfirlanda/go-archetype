package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/dto/response/task"

	"github.com/gofiber/fiber/v2"
)

// GetByPublicID godoc
// @Summary      Get task by ID
// @Description  Retrieve a single task by its ID
// @Tags         tasks
// @Produce      json
// @Security     JWTAuth
// @Security     ApiKeyAuth
// @Param        public_id   path     string  true  "Entity Public ID"
// @Success      200  {object} response.Success{data=taskresp.Detail}
// @Failure      400  {object} response.Error
// @Failure      404  {object} response.Error
// @Router       /api/v1/tasks/{public_id} [get]
func (h *Handler) GetByPublicID(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	publicID := c.Params("public_id")
	if publicID == "" {
		log.Error("task ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("task publicID is required", rid))
	}

	taskEntity, err := h.taskService.GetByPublicID(c.UserContext(), publicID)
	if err != nil {
		return err
	}

	dto := taskresp.ToDetail(taskEntity)
	return c.JSON(response.OK(dto, rid))
}
