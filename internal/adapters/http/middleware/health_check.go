package middleware

import (
	"go-archetype/internal/infrastructure/logging"
	portout "go-archetype/internal/ports/output"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)
import "github.com/gofiber/fiber/v2/middleware/healthcheck"

func HealthCheck(logger *logrus.Entry, dbPinger portout.DBPinger) fiber.Handler {
	log := logging.WithComponent(logger, "middleware.HealthCheck")

	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			// TODO: in future redis, kafka/rabbitmq, or external service
			var errDBPing error
			if dbPinger != nil {
				errDBPing = dbPinger.Ping(c.UserContext())
				if errDBPing != nil {
					log.WithError(errDBPing).Error("failed to ping db")
				}
			}

			return errDBPing == nil
		},
		ReadinessEndpoint: "/ready",
	})
}
