package fiberhttp

import (
	"fmt"
	"go-archetype/internal/adapter/http/fiber/middleware"
	"go-archetype/internal/config"
	"go-archetype/internal/domain/health"
	"go-archetype/internal/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
)

type Dependencies struct {
	APIKeyMiddleware fiber.Handler
	JWTMiddleware    fiber.Handler
	DBPinger         health.DBPinger
}

func StartServer(cfg *config.Config, logger *logrus.Entry, dependencies Dependencies) error {
	log := logging.WithComponent(logger, "http.server")
	app := fiber.New(fiber.Config{
		AppName:      cfg.AppName,
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
	apiKeyMiddleware := middleware.AuthAPIKey(log, cfg.Services.General.APIKey)
	dependencies.APIKeyMiddleware = apiKeyMiddleware
	jwtMiddleware := middleware.AuthJWT(log, cfg.JWT.Secret)
	dependencies.JWTMiddleware = jwtMiddleware

	// Register routes
	RegisterRoutes(app, cfg, log, dependencies)

	log.Infof("Starting HTTP server on port %d", cfg.Http.Port)
	return app.Listen(fmt.Sprintf(":%d", cfg.Http.Port))
}
