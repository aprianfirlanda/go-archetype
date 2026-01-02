package taskhandler

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/request/task"
	"go-archetype/internal/adapters/http/dto/response"
	taskresp "go-archetype/internal/adapters/http/dto/response/task"
	"go-archetype/internal/adapters/http/validation"

	"github.com/gofiber/fiber/v2"
)

// List godoc
// @Summary      List tasks
// @Description  List tasks with pagination and filters
// @Tags         tasks
// @Produce      json
// @Security     JWTAuth
// @Param        page     query int     false "Page number"
// @Param        limit    query int     false "Page size"
// @Param        search   query string  false "Search keyword"
// @Param        status   query string  false "Entity status"
// @Param        priority query int     false "Entity priority"
// @Success      200 {object} response.Success{data=[]taskresp.ListItem, meta=response.PaginationMeta}
// @Failure      400 {object} response.Error{errors=taskresp.ListValidateError}
// @Router       /v1/api/tasks [get]
func (h *Handler) List(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	rid := httpctx.GetRequestID(c)

	var q taskreq.List
	if err := c.QueryParser(&q); err != nil {
		log.WithError(err).Error("failed to parse query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse query params", rid))
	}

	fieldErrors, err := validation.ValidateStruct(q)
	if err != nil {
		log.WithError(err).Error("failed to validate query params")
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate query params", rid))
	}
	if fieldErrors != nil {
		log.WithError(err).Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	q.Normalize()

	filter, err := q.ToListFilter()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(response.FailMessage(err.Error(), rid))
	}

	tasks, total, err := h.taskService.List(c.Context(), filter)
	if err != nil {
		return err
	}

	dto := taskresp.ToList(tasks)
	meta := response.NewPaginationMeta(
		filter.Page,
		filter.Limit,
		total,
	)

	return c.Status(fiber.StatusOK).
		JSON(response.OKPaginate(dto, meta, rid))
}
