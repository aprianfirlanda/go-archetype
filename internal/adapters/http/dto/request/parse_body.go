package httpreq

import (
	"go-archetype/internal/adapters/http/dto/response"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func ParseBody[T any](c *fiber.Ctx, log *logrus.Entry, rid string) (T, error) {
	var req T
	if err := c.BodyParser(&req); err != nil {
		log.WithError(err).Error("failed to parse request body")
		return req, c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to parse request body", rid))
	}

	if err := validate(req, c, log, rid, "request body"); err != nil {
		return req, err
	}

	return req, nil
}
