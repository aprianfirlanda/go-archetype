package fiberhttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func RegisterRoutes(app *fiber.App, logger *logrus.Logger, dependencies Dependencies) {

	app.Get("/protected-by-api-key", dependencies.APIKeyMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("Hello from protected route by API key!")
	})
	app.Get("/panic", func(c *fiber.Ctx) error {
		logger.Info("About to panic with nil pointer")

		var x *int
		// this will cause: panic: runtime error: invalid memory address or nil pointer dereference
		_ = *x

		return nil
	})
}
