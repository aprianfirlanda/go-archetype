package fiberhttp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverPanic "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/config"
)

func StartServer(appName string, httpConfig config.Http, logger *logrus.Logger) error {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	// Global middlewares
	app.Use(requestid.New())
	app.Use(recoverPanic.New())
	app.Use(cors.New())

	// Register routes
	RegisterRoutes(app, logger)

	return app.Listen(fmt.Sprintf(":%d", httpConfig.Port))
}
