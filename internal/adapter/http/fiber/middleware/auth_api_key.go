package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"github.com/gofiber/fiber/v2"
)
import "github.com/gofiber/fiber/v2/middleware/keyauth"

func AuthAPIKey(apiKey string) fiber.Handler {
	hashedAPIKey := sha256.Sum256([]byte(apiKey))

	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, incomingAPIKey string) (bool, error) {
			incomingHash := sha256.Sum256([]byte(incomingAPIKey))

			if subtle.ConstantTimeCompare(hashedAPIKey[:], incomingHash[:]) == 1 {
				return true, nil
			}
			return false, keyauth.ErrMissingOrMalformedAPIKey
		},
	})
}
