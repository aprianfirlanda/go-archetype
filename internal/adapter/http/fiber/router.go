package fiberhttp

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/adapter/http/fiber/middleware"
	"go-archetype/internal/config"
	"go-archetype/internal/domain/auth"
	"time"
)

func RegisterRoutes(app *fiber.App, cfg *config.Config, logger *logrus.Logger, dependencies Dependencies) {
	app.Get("/protected-by-api-key", dependencies.APIKeyMiddleware, func(c *fiber.Ctx) error {
		log := middleware.RequestLogger(c, logger)
		log.Info("Hello from protected route by API key!")
		return c.SendString("Hello from protected route by API key!")
	})
	app.Get("/generate-token", func(c *fiber.Ctx) error {
		log := middleware.RequestLogger(c, logger)

		claims := auth.CustomClaims{}
		claims.Roles = []string{"admin"}
		claims.RegisteredClaims = jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Subject:   "user-123",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    cfg.AppName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := token.SignedString([]byte(cfg.JWT.Secret))
		if err != nil {
			log.WithError(err).Error("failed to generate jwt")
			return c.Status(500).JSON(fiber.Map{"error": "internal error"})
		}

		log.WithField("user_id", claims.Subject).Info("login success")

		return c.JSON(fiber.Map{"token": signedToken})
	})
	app.Get("/protected-by-jwt", dependencies.JWTMiddleware, func(c *fiber.Ctx) error {
		return c.SendString("Hello from protected route by JWT!")
	})
	app.Get("/panic", func(c *fiber.Ctx) error {
		logger.Info("About to panic with nil pointer")

		var x *int
		// this will cause: panic: runtime error: invalid memory address or nil pointer dereference
		_ = *x

		return nil
	})
}
