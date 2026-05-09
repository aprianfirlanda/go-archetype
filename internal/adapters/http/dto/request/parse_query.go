package httpreq

import (
	"go-archetype/internal/adapters/http/dto/response"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ParseQuery[T any](c *fiber.Ctx, log *logrus.Entry, rid string) (T, error) {
	var req T
	if err := c.QueryParser(&req); err != nil {
		log.WithError(err).Error("failed to parse query params")
		return req, c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse query params", rid))
	}

	if err := validate(req, c, log, rid, "query params"); err != nil {
		return req, err
	}

	return req, nil
}
