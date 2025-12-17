package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RequestIDLoggerContext(logger *logrus.Entry) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rid := c.Locals("requestid").(string)
		entry := logger.WithField("request_id", rid)
		c.Locals("logger", entry)

		return c.Next()
	}
}

func RequestLogger(c *fiber.Ctx, fallback *logrus.Entry) *logrus.Entry {
	if v := c.Locals("logger"); v != nil {
		if entry, ok := v.(*logrus.Entry); ok {
			if fallback == nil || len(fallback.Data) == 0 {
				return entry
			}

			return entry.WithFields(fallback.Data)
		}
	}
	return fallback
}

func GetRequestID(c *fiber.Ctx) string {
	if v := c.Locals("requestid"); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
