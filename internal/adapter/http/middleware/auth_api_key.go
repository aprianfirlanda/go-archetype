package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	httpctx "go-archetype/internal/adapter/http/context"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func validateAPIKey(c *fiber.Ctx, apiKey string) bool {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return false
	}

	// Remove "ApiKey " prefix
	incomingAPIKey := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "ApiKey " {
		incomingAPIKey = authHeader[7:]
	} else {
		return false // Wrong format
	}

	hashedAPIKey := sha256.Sum256([]byte(apiKey))
	incomingHash := sha256.Sum256([]byte(incomingAPIKey))

	return subtle.ConstantTimeCompare(hashedAPIKey[:], incomingHash[:]) == 1
}

// AuthAPIKey returns a handler for use with AnyAuth
func AuthAPIKey(logger *logrus.Entry, apiKey string) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "http.middleware.AuthAPIKey")

	return func(c *fiber.Ctx) error {
		log := httpctx.Get(c, logWithComponent)

		if validateAPIKey(c, apiKey) {
			log.Info("API key validated successfully")
			return nil // Success
		}

		return fiber.NewError(fiber.StatusUnauthorized, "invalid API key")
	}
}
