package middleware

import (
	"go-archetype/internal/adapters/http/context"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestIDContext(base *logrus.Entry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := httpctx.EnrichRequestID(base, c)

		httpctx.Set(c, log)
		return c.Next()
	}
}
