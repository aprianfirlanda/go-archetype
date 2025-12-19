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

	api := app.Group("/api")

	// Demo API
	demoHandler := handler.NewDemoHandler(log, deps.Config)
	demo := api.Group("/demo")
	demo.Get("/protected-by-api-key", apiKeyMiddleware, demoHandler.ProtectedByAPIKey)
	demo.Get("/generate-token", demoHandler.GenerateToken)
	demo.Get("/protected-by-jwt", jwtMiddleware, demoHandler.ProtectedByJWT)
	demo.Get("/panic", demoHandler.Panic)

	// Task API
	taskHandler := handler.NewTaskHandler(log)
	task := api.Group("/tasks")
	task.Post("/", jwtMiddleware, taskHandler.Create)
	task.Get("/", middleware.AnyAuth(apiKeyMiddleware, apiKeyMiddleware), taskHandler.List)
	task.Get("/:id", jwtMiddleware, taskHandler.GetByID)
	task.Put("/:id", jwtMiddleware, taskHandler.Update)
	task.Patch("/:id/status", jwtMiddleware, taskHandler.UpdateStatus)
	task.Patch("/status", jwtMiddleware, taskHandler.BulkUpdateStatus)
	task.Delete("/:id", jwtMiddleware, taskHandler.Delete)
	task.Delete("/", jwtMiddleware, taskHandler.BulkDelete)
}
