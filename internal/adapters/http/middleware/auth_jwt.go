package middleware

import (
	"fmt"
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/domain/auth"
	"go-archetype/internal/infrastructure/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

// validateJWT checks if the JWT is valid
func validateJWT(c *fiber.Ctx, jwtSecret string) bool {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return false
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	jwtSecretBytes := []byte(jwtSecret)
	token, err := jwt.ParseWithClaims(tokenString, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecretBytes, nil
	})

	if err != nil {
		return false
	}

	claims, ok := token.Claims.(*auth.CustomClaims)
	if !ok || !token.Valid {
		return false
	}

	// Store claims in context
	c.Locals("user", claims)
	return true
}

// AuthJWT returns a handler for use with AnyAuth
func AuthJWT(logger *logrus.Entry, jwtSecret string) fiber.Handler {
	logWithComponent := logging.WithComponent(logger, "http.middleware.AuthJWT")

	return func(c *fiber.Ctx) error {
		log := httpctx.Get(c, logWithComponent)

		if validateJWT(c, jwtSecret) {
			log.Info("JWT validated successfully")
			return nil // Success
		}

		return fiber.NewError(fiber.StatusUnauthorized, "invalid JWT")
	}
}
