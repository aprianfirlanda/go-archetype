package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
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
// @Router       /api/v1/tasks/status [patch]
func (h *Handler) BulkUpdateStatus(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	req, err := httpreq.ParseBody[taskreq.BulkUpdateStatus](c, log, rid)
	if err != nil {
		return err
	}

	cmd := taskcmd.BulkUpdateStatus{
		PublicIDs: req.IDs,
		Status:    task.Status(req.Status),
	}
	res, err := h.taskService.BulkUpdateStatus(c.UserContext(), cmd)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.OK(res, rid))
}
