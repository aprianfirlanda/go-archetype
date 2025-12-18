package router

import (
	"go-archetype/internal/adapter/http/handler"
	"go-archetype/internal/adapter/http/middleware"
	"go-archetype/internal/bootstrap"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, deps bootstrap.HttpApp) {
	log := logging.WithComponent(deps.Log, "http.router")

	// Setup Auth Middlewares
	apiKeyMiddleware := middleware.AuthAPIKey(log, deps.Config.Services.General.APIKey)
	jwtMiddleware := middleware.AuthJWT(log, deps.Config.JWT.Secret)

	// Setup Handler
	demoHandler := handler.NewDemoHandler(log, deps.Config)

	app.Get("/protected-by-api-key", apiKeyMiddleware, demoHandler.ProtectedByAPIKey)
	app.Get("/generate-token", demoHandler.GenerateToken)
	app.Get("/protected-by-jwt", jwtMiddleware, demoHandler.ProtectedByJWT)
	app.Get("/panic", demoHandler.Panic)
}
