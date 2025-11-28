package middleware

import "github.com/gofiber/fiber/v2"
import fiberhealthcheck "github.com/gofiber/fiber/v2/middleware/healthcheck"

func HealthCheck() fiber.Handler {
	return fiberhealthcheck.New(fiberhealthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			// TODO: in future check db, redis, kafka/rabbitmq, or external service
			// TODO: put on log the app that not ready
			return true
		},
		ReadinessEndpoint: "/ready",
	})
}
