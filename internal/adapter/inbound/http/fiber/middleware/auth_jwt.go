package middleware

import (
	"fmt"
	"go-archetype/internal/adapter/inbound/http/fiber/response"
	"go-archetype/internal/domain/auth"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func AuthJWT(logger *logrus.Entry, jwtSecret string) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "middleware.AuthJWT")
	jwtSecretBytes := []byte(jwtSecret)

	return keyauth.New(keyauth.Config{
		Validator: func(c *fiber.Ctx, tokenString string) (bool, error) {
			log := RequestLogger(c, logWithComponent)
			// Parse & validate token
			token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				// Ensure the signing method is what you expect
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					log.Warn("Unexpected signing method: ", t.Header["alg"])
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return jwtSecretBytes, nil
			})
			if err != nil {
				log.WithError(err).Warn("failed to parse JWT")
				return false, err
			}

			// Token must be valid and have claims
			claims, ok := token.Claims.(*auth.CustomClaims)
			if !ok || !token.Valid {
				log.Warn("invalid token")
				return false, fmt.Errorf("invalid token")
			}

			// You can do extra checks here if needed:
			// - check Role
			// - check subject / audience
			// - check custom fields

			// Store claims in context so handlers can use it
			c.Locals("user", claims)

			log.Info("jwt validated successfully")

			return true, nil
		},
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			status := fiber.StatusUnauthorized

			// get requestId from fiber/middleware/requestid
			rid := GetRequestID(c)

			resp := response.ErrorResponse{
				Message:   "invalid or missing JWT",
				RequestID: rid,
			}

			c.Type("json", "utf-8")
			return c.Status(status).JSON(resp)
		},
	})
}
