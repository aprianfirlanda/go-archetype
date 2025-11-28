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

		// Get request ID from locals (must match your requestid middleware ContextKey)
		requestID, _ := c.Locals("requestid").(string)

		status := c.Response().StatusCode()
		method := c.Method()
		path := c.OriginalURL()

		// Build structured log entry
		entry := logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"status":     status,
			"method":     method,
			"path":       path,
			"latency_ms": float64(latency.Microseconds()) / 1000.0,
			"ip":         c.IP(),
		})

		if err != nil {
			entry = entry.WithError(err)
		}

		// Log level based on status
		switch {
		case status >= 500:
			entry.Error("HTTP request completed")
		case status >= 400:
			entry.Warn("HTTP request completed")
		default:
			entry.Info("HTTP request completed")
		}

		return err
	}
}
