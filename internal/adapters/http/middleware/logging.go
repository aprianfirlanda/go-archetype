package middleware

import (
	"errors"
	"time"

	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func Logging(base *logrus.Entry) fiber.Handler {
	componentLog := logging.WithComponent(base, "http.middleware.logging")

	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		log := httpctx.Get(c, componentLog)
		log = httpctx.EnrichMeta(log, c, latency.Milliseconds())
		log = httpctx.EnrichUserInfo(log, c)

		status := resolveStatus(c, err)

		if err != nil {
			log = log.WithError(err)
		}

		log = log.WithField("status", status)

		switch {
		case status >= 500:
			log.Error("http request completed")
		case status >= 400:
			log.Warn("http request completed")
		default:
			log.Info("http request completed")
		}

		return err
	}
}

func resolveStatus(c *fiber.Ctx, err error) int {
	if err == nil {
		return c.Response().StatusCode()
	}

	// Fiber-native error
	var fe *fiber.Error
	if errors.As(err, &fe) {
		return fe.Code
	}

	// Unknown error
	return fiber.StatusInternalServerError
}
