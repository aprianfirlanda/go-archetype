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
			log := RequestLogger(c, logger)
			log.WithFields(logrus.Fields{
				"error":      e,
				"stacktrace": string(debug.Stack()),
			}).Error("panic recovered")
		},
	})
}
