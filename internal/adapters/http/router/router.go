package router

import (
	"go-archetype/internal/adapters/http/handler/demo"
	"go-archetype/internal/adapters/http/handler/task"
	"go-archetype/internal/adapters/http/middleware"
	"go-archetype/internal/bootstrap"
	"go-archetype/internal/infrastructure/logging"

	_ "go-archetype/internal/adapters/http/docs"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func RegisterRoutes(app *fiber.App, deps bootstrap.HttpApp) {
	log := logging.WithComponent(deps.Log, "http.router")

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Setup Auth Middlewares
	apiKeyMiddleware := middleware.AuthAPIKey(log, deps.Config.Services.General.APIKey)
	jwtMiddleware := middleware.AuthJWT(log, deps.Config.JWT.Secret)

	api := app.Group("/api")

	// Demo API
	demoHandler := demohandler.NewDemoHandler(log, deps.Config)
	demoV1 := api.Group("/v1/demo")
	demoV1.Get("/protected-by-api-key", apiKeyMiddleware, demoHandler.ProtectedByAPIKey)
	demoV1.Get("/generate-token", demoHandler.GenerateToken)
	demoV1.Get("/protected-by-jwt", jwtMiddleware, demoHandler.ProtectedByJWT)
	demoV1.Get("/panic", demoHandler.Panic)

	// Task API
	taskHandler := taskhandler.NewHandler(log, deps.TaskService)
	taskV1 := api.Group("/v1/tasks")
	taskV1.Post("/", jwtMiddleware, taskHandler.Create)
	taskV1.Get("/", middleware.AnyAuth(apiKeyMiddleware, jwtMiddleware), taskHandler.List)
	taskV1.Get("/:public_id", jwtMiddleware, taskHandler.GetByPublicID)
	taskV1.Put("/:public_id", jwtMiddleware, taskHandler.Update)
	taskV1.Patch("/:public_id/status", jwtMiddleware, taskHandler.UpdateStatus)
	taskV1.Patch("/status", jwtMiddleware, taskHandler.BulkUpdateStatus)
	taskV1.Delete("/:public_id", jwtMiddleware, taskHandler.DeletePublicID)
	taskV1.Delete("/", jwtMiddleware, taskHandler.BulkDelete)
}
