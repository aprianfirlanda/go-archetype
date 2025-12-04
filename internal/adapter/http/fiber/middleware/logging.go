package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

func Logging(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Continue to the next middleware/handler
		err := c.Next()

		latency := time.Since(start)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()

		// Build structured log entry
		log := RequestLogger(c, logger)
		log.WithFields(logrus.Fields{
			"status":     status,
			"method":     method,
			"path":       path,
			"latency_ms": float64(latency.Microseconds()) / 1000.0,
			"ip":         c.IP(),
		})

		if err != nil {
			log.WithError(err)
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
