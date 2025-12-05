package fiberhttp

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/adapter/http/fiber/middleware"
	"go-archetype/internal/config"
)

type Dependencies struct {
	APIKeyMiddleware fiber.Handler
	JWTMiddleware    fiber.Handler
}

func StartServer(cfg *config.Config, logger *logrus.Logger, dependencies Dependencies) error {
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
	})

	// Global middlewares
	// 0. Health Check, live: is the application up, ready: is the application ready to accept traffic
	app.Use(middleware.HealthCheck())
	app.Get("/metrics", monitor.New())
	// 1. Generate request ID first so everyone can use it
	app.Use(requestid.New())
	app.Use(middleware.RequestIDLoggerContext(logger))
	// 2. Logging wraps everything below (including recover + cors + handlers)
	app.Use(middleware.Logging(logger))
	// 3. Recover from panic so we don't crash the server
	app.Use(middleware.Recover(logger))
	// 4. CORS – mostly for browser / frontend
	app.Use(cors.New())

	// Auth Middleware
	apiKeyMiddleware := middleware.AuthAPIKey(logger, cfg.Services.General.APIKey)
	dependencies.APIKeyMiddleware = apiKeyMiddleware
	jwtMiddleware := middleware.AuthJWT(logger, cfg.JWT.Secret)
	dependencies.JWTMiddleware = jwtMiddleware

	// Register routes
	RegisterRoutes(app, cfg, logger, dependencies)

	logger.Infof("Starting HTTP server on port %d", cfg.Http.Port)
	return app.Listen(fmt.Sprintf(":%d", cfg.Http.Port))
}
