package fiberhttp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/adapter/http/fiber/middleware"
	"go-archetype/internal/config"
)

type Dependencies struct{}

func StartServer(appName string, httpConfig config.Http, logger *logrus.Logger, dependencies Dependencies) error {
	app := fiber.New(fiber.Config{
		AppName: appName,
	})

	// Global middlewares
	// 0. Health Check, live: is the application up, ready: is the application ready to accept traffic
	app.Use(middleware.HealthCheck())
	// 1. Generate request ID first so everyone can use it
	app.Use(requestid.New())
	// 2. Logging wraps everything below (including recover + cors + handlers)
	app.Use(middleware.Logging(logger))
	// 3. Recover from panic so we don't crash the server
	app.Use(middleware.Recover(logger))
	// 4. CORS – mostly for browser / frontend
	app.Use(cors.New())

	// Register routes
	RegisterRoutes(app, logger, dependencies)

	logger.Infof("Starting HTTP server on port %d", httpConfig.Port)
	return app.Listen(fmt.Sprintf(":%d", httpConfig.Port))
}
