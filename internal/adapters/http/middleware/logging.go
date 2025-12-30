package middleware

import (
	"time"

	httpctx "go-archetype/internal/adapters/http/context"
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

		if err != nil {
			log = log.WithError(err)
		}

		switch {
		case c.Response().StatusCode() >= 500:
			log.Error("http request completed")
		case c.Response().StatusCode() >= 400:
			log.Warn("http request completed")
		default:
			log.Info("http request completed")
		}

		return err
	}
}
