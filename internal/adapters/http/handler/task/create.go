package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/application/task/command"

	"github.com/gofiber/fiber/v2"
)

// Create godoc
// @Summary      Create a task
// @Description  Create a new task
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     JWTAuth
// @Param        request body taskreq.Create true "Create Entity Request"
// @Success      201 {object} response.Success{data=response.IDResponse}
// @Failure      400 {object} response.Error{errors=taskresp.CreateValidateError}
// @Failure      500 {object} response.Error
// @Router       /api/v1/tasks [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	req, err := httpreq.ParseBody[taskreq.Create](c, log, rid)
	if err != nil {
		return err
	}

	cmd := taskcmd.Create{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
		Tags:        req.Tags,
	}
	publicID, err := h.taskService.Create(c.UserContext(), cmd)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.OK(response.IDResponse{ID: publicID}, rid))
}
