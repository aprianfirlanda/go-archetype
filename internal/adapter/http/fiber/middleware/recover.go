package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func Recover(logger *logrus.Logger) fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// get request id
			rid, _ := c.Locals("requestid").(string)

			logger.WithFields(logrus.Fields{
				"request_id": rid,
				"error":      e,
				"method":     c.Method(),
				"path":       c.OriginalURL(),
				"stacktrace": string(debug.Stack()),
			}).Error("panic recovered")
		},
	})
}
