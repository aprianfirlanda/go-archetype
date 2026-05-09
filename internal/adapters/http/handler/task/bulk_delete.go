package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
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
// @Router       /api/v1/tasks [delete]
func (h *Handler) BulkDelete(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	req, err := httpreq.ParseBody[taskreq.BulkDelete](c, log, rid)
	if err != nil {
		return err
	}

	cmd := taskcmd.BulkDelete{
		PublicIDs: req.IDs,
	}
	res, err := h.taskService.BulkDelete(c.UserContext(), cmd)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.OK(res, rid))
}
