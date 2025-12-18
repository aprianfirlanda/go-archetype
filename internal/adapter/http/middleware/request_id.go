package middleware

import (
	httpctx "go-archetype/internal/adapter/http/context"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestIDContext(base *logrus.Entry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := logging.WithComponent(base, "http.middleware.request")
		log = httpctx.EnrichRequestID(log, c)

		httpctx.Set(c, log)
		return c.Next()
	}
}
