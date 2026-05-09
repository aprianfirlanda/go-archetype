package httpreq

import (
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/adapters/http/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func validate(req any, c *fiber.Ctx, log *logrus.Entry, rid string, source string) error {
	fieldErrors, err := validation.ValidateStruct(req)
	if err != nil {
		log.WithError(err).Errorf("failed to validate %s", source)
		return c.Status(fiber.StatusBadRequest).JSON(response.FailMessage("failed to validate "+source, rid))
	}
	if fieldErrors != nil {
		log.Error("validation failed")
		return c.Status(fiber.StatusBadRequest).JSON(response.Fail("validation failed", fieldErrors, rid))
	}

	return nil
}
