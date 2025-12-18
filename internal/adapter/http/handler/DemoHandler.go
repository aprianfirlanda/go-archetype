package handler

import (
	httpctx "go-archetype/internal/adapter/http/context"
	"go-archetype/internal/domain/auth"
	"go-archetype/internal/infrastructure/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type DemoHandler struct {
	log *logrus.Entry
	cfg *config.Config
}

func NewDemoHandler(log *logrus.Entry, cfg *config.Config) *DemoHandler {
	return &DemoHandler{
		log: log,
		cfg: cfg,
	}
}

func (h *DemoHandler) ProtectedByAPIKey(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	log.Info("Hello from protected route by API key!")
	return c.SendString("Hello from protected route by API key!")
}

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

func (h *DemoHandler) ProtectedByJWT(c *fiber.Ctx) error {
	return c.SendString("Hello from protected route by JWT!")
}

func (h *DemoHandler) Panic(_ *fiber.Ctx) error {
	h.log.Warn("About to panic with nil pointer")

	var x *int
	_ = *x
	return nil
}
