package demohandler

import (
	httpctx "go-archetype/internal/adapters/http/context"
	"go-archetype/internal/domain/auth"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (h *DemoHandler) GenerateToken(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)

	claims := auth.CustomClaims{
		Roles: []string{"admin"},
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   "user-123",
			Issuer:    h.cfg.AppName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(h.cfg.JWT.Secret))
	if err != nil {
		log.WithError(err).Error("failed to generate jwt")
		return fiber.ErrInternalServerError
	}

	log.WithField("user_id", claims.Subject).Info("login success")
	return c.JSON(fiber.Map{"token": signed})
}
