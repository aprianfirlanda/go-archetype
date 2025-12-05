package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"go-archetype/internal/domain/auth"
)
import "github.com/gofiber/fiber/v2/middleware/keyauth"

func AuthJWT(logger *logrus.Logger, jwtSecret string) fiber.Handler {
	jwtSecretBytes := []byte(jwtSecret)

	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, tokenString string) (bool, error) {
			log := RequestLogger(c, logger)
			// Parse & validate token
			token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				// Ensure the signing method is what you expect
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Error("Unexpected signing method: ", t.Header["alg"])
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return jwtSecretBytes, nil
			})
			if err != nil {
				log.WithError(err).Error("failed to parse JWT")
				return false, err
			}

			// Token must be valid and have claims
			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok || !token.Valid {
				log.Error("invalid token")
				return false, fmt.Errorf("invalid token")
			}

			// You can do extra checks here if needed:
			// - check Role
			// - check subject / audience
			// - check custom fields

			// Store claims in context so handlers can use it
			c.Locals("user", claims)

			// Enrich Logger: add user detail on the SAME log entry
			updatedLog := log.WithFields(logrus.Fields{
				"subject": claims.Subject,
				"roles":   claims.Roles,
			})

			// overwrite contextual logger for downstream
			c.Locals("logger", updatedLog)

			updatedLog.Info("jwt validated successfully")

			return true, nil
		},
	})
}
