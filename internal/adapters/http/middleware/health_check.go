package middleware

import (
	"go-archetype/internal/infrastructure/logging"
	portin "go-archetype/internal/ports/input"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)
import "github.com/gofiber/fiber/v2/middleware/healthcheck"

func HealthCheck(logger *logrus.Entry, svc portin.HealthService) fiber.Handler {
	log := logging.WithComponent(logger, "middleware.HealthCheck")

	return healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return svc.Liveness(c.UserContext())
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			if err := svc.Readiness(c.UserContext()); err != nil {
				log.WithError(err).Error("readiness check failed")
				return false
			}
			return true
		},
		ReadinessEndpoint: "/ready",
	})
}
