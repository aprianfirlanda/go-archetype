package middleware

import (
	"go-archetype/internal/infrastructure/logging"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

func Recover(logger *logrus.Entry) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "middleware.Recover")
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			log := RequestLogger(c, logWithComponent)
			log.WithFields(logging.Fields(map[string]any{
				"error":      e,
				"stacktrace": string(debug.Stack()),
			})).Error("panic recovered")
		},
	})
}
