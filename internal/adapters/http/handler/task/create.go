package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"
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

	var req taskreq.Create
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
