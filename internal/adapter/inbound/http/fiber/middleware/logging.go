package middleware

import (
	"go-archetype/internal/domain/auth"
	"go-archetype/internal/infrastructure/logging"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logging(logger *logrus.Entry) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "middleware.Logging")
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to the next middleware/handler
		err := c.Next()

		latency := time.Since(start)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()

		// Build structured log entry
		log := RequestLogger(c, logWithComponent).WithFields(logrus.Fields{
			"status":     status,
			"method":     method,
			"path":       path,
			"latency_ms": float64(latency.Microseconds()) / 1000.0,
			"ip":         c.IP(),
		})

		if v := c.Locals("user"); v != nil {
			claims, ok := v.(*auth.CustomClaims)
			if ok {
				log = log.WithFields(logrus.Fields{
					"user_id": claims.Subject,
					"roles":   claims.Roles,
				})
			} else {
				log.Error("failed to cast user claims to CustomClaims")
			}
		}

		if err != nil {
			log = log.WithError(err)
		}

		// Log level based on status
		switch {
		case status >= 500:
			log.Error("HTTP request completed")
		case status >= 400:
			log.Warn("HTTP request completed")
		default:
			log.Info("HTTP request completed")
		}

		return err
	}
}
