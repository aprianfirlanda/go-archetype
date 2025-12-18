package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	httpctx "go-archetype/internal/adapter/http/context"
	"go-archetype/internal/adapter/http/dto/response"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)
import "github.com/gofiber/fiber/v2/middleware/keyauth"

func AuthAPIKey(logger *logrus.Entry, apiKey string) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "http.middleware.AuthAPIKey")
	hashedAPIKey := sha256.Sum256([]byte(apiKey))

	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, incomingAPIKey string) (bool, error) {
			incomingHash := sha256.Sum256([]byte(incomingAPIKey))

			if subtle.ConstantTimeCompare(hashedAPIKey[:], incomingHash[:]) == 1 {
				return true, nil
			}

			log := httpctx.Get(c, logWithComponent)
			log.WithFields(logging.Field("reason", "invalid API key")).Warn("unauthorized request")
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusUnauthorized

			// get requestId from fiber/middleware/requestid
			rid := httpctx.GetRequestID(c)

			resp := response.ErrorResponse{
				Message:   "invalid or missing API key",
				RequestID: rid,
			}

			c.Type("json", "utf-8")
			return c.Status(status).JSON(resp)
		},
	})
}
