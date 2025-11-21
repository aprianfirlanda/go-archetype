package fiberhttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RegisterRoutes(app *fiber.App, logger *logrus.Logger, dependencies Dependencies) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		logger.WithFields(logrus.Fields{
			"request_id": c.Locals("requestid"),
		}).Info("Health check")
		return c.JSON(fiber.Map{"status": "ok"})
	})
}
