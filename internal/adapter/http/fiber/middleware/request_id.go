package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestIDLoggerContext(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Locals("requestid").(string)
		entry := logger.WithField("request_id", rid)
		c.Locals("logger", entry)

		return c.Next()
	}
}

func RequestLogger(c *fiber.Ctx, fallback *logrus.Logger) *logrus.Entry {
	if v := c.Locals("logger"); v != nil {
		if entry, ok := v.(*logrus.Entry); ok {
			return entry
		}
	}
	return logrus.NewEntry(fallback)
}
