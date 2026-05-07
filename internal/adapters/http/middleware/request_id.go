package middleware

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestIDContext(base *logrus.Entry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := httpctx.EnrichRequestID(base, c)
		rid := httpctx.GetRequestID(c)

		httpctx.Set(c, log)
		ctx := logging.WithLogger(c.UserContext(), log)
		ctx = logging.WithRequestID(ctx, rid)
		c.SetUserContext(ctx)
		return c.Next()
	}
}
