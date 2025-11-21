package fiberhttp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverPanic "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/adapter/http/fiber/middleware"
	"go-archetype/internal/config"
)

func StartServer(appName string, httpConfig config.Http, logger *logrus.Logger) error {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	// Global middlewares
	// 1. Generate request ID first so everyone can use it
	app.Use(requestid.New())
	// 2. Logging wraps everything below (including recover + cors + handlers)
	app.Use(middleware.Logging(logger))
	// 3. Recover from panic so we don't crash the server
	app.Use(recoverPanic.New())
	// 4. CORS – mostly for browser / frontend
	app.Use(cors.New())

	// Register routes
	RegisterRoutes(app, logger)

	return app.Listen(fmt.Sprintf(":%d", httpConfig.Port))
}
