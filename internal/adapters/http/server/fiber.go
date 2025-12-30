package server

import (
	"fmt"
	"go-archetype/internal/adapters/http/middleware"
	"go-archetype/internal/adapters/http/router"
	"go-archetype/internal/bootstrap"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func StartServer(deps bootstrap.HttpApp) error {
	log := logging.WithComponent(deps.Log, "http.server")
	app := fiber.New(fiber.Config{
		AppName:      deps.Config.AppName,
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Global middlewares
	// 0. Health Check, live: is the application up, ready: is the application ready to accept traffic
	app.Use(middleware.HealthCheck(log, deps.DBPinger))
	app.Get("/metrics", monitor.New())
	// 1. Generate request ID first so everyone can use it
	app.Use(requestid.New())
	app.Use(middleware.RequestIDContext(log))
	// 2. Logging wraps everything below (including recover + cors + handlers)
	app.Use(middleware.Logging(log))
	// 3. Recover from panic so we don't crash the server
	app.Use(middleware.Recover(log))
	// 4. CORS – mostly for browser / frontend
	app.Use(cors.New())

	// Register routes
	router.RegisterRoutes(app, deps)

	log.Infof("Starting HTTP server on port %d", deps.Config.Http.Port)
	return app.Listen(fmt.Sprintf(":%d", deps.Config.Http.Port))
}
