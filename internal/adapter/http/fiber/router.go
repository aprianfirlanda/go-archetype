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
	app.Get("/panic", func(c *fiber.Ctx) error {
		logger.Info("About to panic with nil pointer")

		var x *int
		// this will cause: panic: runtime error: invalid memory address or nil pointer dereference
		_ = *x

		return nil
	})
}
