package middleware

import (
	"github.com/gofiber/fiber/v2"
	fiberrecover "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func Recover(logger *logrus.Logger) fiber.Handler {
	return fiberrecover.New(fiberrecover.Config{
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
