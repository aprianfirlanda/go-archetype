package middleware

import (
	"context"
	"crypto/tls"
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/infrastructure/config"
	"go-archetype/internal/infrastructure/logging"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func AuthKeycloak(logger *logrus.Entry, cfg config.Keycloak) fiber.Handler {

	logWithComponent := logging.WithComponent(logger, "http.middleware.AuthKeycloak")

	// Init OIDC provider once
	insecureClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.InsecureSkipVerify,
			},
		},
	}
	ctx := oidc.ClientContext(context.Background(), insecureClient)
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		panic("failed to init keycloak oidc provider: " + err.Error())
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	return func(c *fiber.Ctx) error {
		log := httpctx.Get(c, logWithComponent)

		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.NewError(fiber.StatusUnauthorized, "missing bearer token")
		}

		rawToken := strings.TrimPrefix(authHeader, "Bearer ")

		accessToken, err := verifier.Verify(c.Context(), rawToken)
		if err != nil {
			log.WithError(err).Warn("invalid keycloak token")
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token")
		}

		var claims map[string]any
		if err := accessToken.Claims(&claims); err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		c.Locals("user", claims)

		log.Info("keycloak JWT validated successfully")
		return c.Next()
	}
}
