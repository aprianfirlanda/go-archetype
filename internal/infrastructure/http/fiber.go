package http

import (
	"fmt"
	"go-archetype/internal/adapter/inbound/http/fiber"
	"go-archetype/internal/adapter/inbound/http/fiber/middleware"
	"go-archetype/internal/bootstrap"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func StartServer(dependencies bootstrap.HttpApp) error {
	log := logging.WithComponent(dependencies.Log, "http.server")
	app := fiber.New(fiber.Config{
		AppName:      dependencies.Config.AppName,
		ErrorHandler: middleware.ErrorHandler(),
	})

	// Global middlewares
	// 0. Health Check, live: is the application up, ready: is the application ready to accept traffic
	app.Use(middleware.HealthCheck(log, dependencies.DBPinger))
	app.Get("/metrics", monitor.New())
	// 1. Generate request ID first so everyone can use it
	app.Use(requestid.New())
	app.Use(middleware.RequestIDLoggerContext(log))
	// 2. Logging wraps everything below (including recover + cors + handlers)
	app.Use(middleware.Logging(log))
	// 3. Recover from panic so we don't crash the server
	app.Use(middleware.Recover(log))
	// 4. CORS – mostly for browser / frontend
	app.Use(cors.New())

	// Auth Middleware
	apiKeyMiddleware := middleware.AuthAPIKey(log, dependencies.Config.Services.General.APIKey)
	dependencies.APIKeyMiddleware = apiKeyMiddleware
	jwtMiddleware := middleware.AuthJWT(log, dependencies.Config.JWT.Secret)
	dependencies.JWTMiddleware = jwtMiddleware

	// Register routes
	fiberhttp.RegisterRoutes(app, dependencies)

	log.Infof("Starting HTTP server on port %d", dependencies.Config.Http.Port)
	return app.Listen(fmt.Sprintf(":%d", dependencies.Config.Http.Port))
}
