package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"go-archetype/internal/adapter/http/fiber/response"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)
import "github.com/gofiber/fiber/v2/middleware/keyauth"

func AuthAPIKey(logger *logrus.Logger, apiKey string) fiber.Handler {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))

	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, incomingAPIKey string) (bool, error) {
			incomingHash := sha256.Sum256([]byte(incomingAPIKey))

			if subtle.ConstantTimeCompare(hashedAPIKey[:], incomingHash[:]) == 1 {
				return true, nil
			}

			log := RequestLogger(c, logger)
			log.WithFields(logrus.Fields{
				"reason": "invalid API key",
			}).Error("unauthorized request")
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusUnauthorized

			// get requestId from fiber/middleware/requestid
			var requestID string
			if v := c.Locals("requestid"); v != nil {
				if id, ok := v.(string); ok {
					requestID = id
				}
			}

			resp := response.ErrorResponse{
				Message:   "invalid or missing API key",
				RequestID: requestID,
			}

			c.Type("json", "utf-8")
			return c.Status(status).JSON(resp)
		},
	})
}
