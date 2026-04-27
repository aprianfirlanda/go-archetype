package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func validateAPIKey(c *fiber.Ctx, apiKeys []string) bool {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return false
	}

	// Expect: "ApiKey xxx"
	if len(authHeader) <= 7 || authHeader[:7] != "ApiKey " {
		return false
	}

	incomingAPIKey := authHeader[7:]
	incomingHash := sha256.Sum256([]byte(incomingAPIKey))

	for _, key := range apiKeys {
		hashed := sha256.Sum256([]byte(key))
		if subtle.ConstantTimeCompare(hashed[:], incomingHash[:]) == 1 {
			return true
		}
	}

	return false
}

// AuthAPIKey returns a handler for use with AnyAuth
func AuthAPIKey(logger *logrus.Entry, apiKeys []string) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "http.middleware.AuthAPIKey")

	return func(c *fiber.Ctx) error {
		log := httpctx.Get(c, logWithComponent)

		if validateAPIKey(c, apiKeys) {
			log.Info("API key validated successfully")
			return c.Next()
		}

		return fiber.NewError(fiber.StatusUnauthorized, "invalid API key")
	}
}
