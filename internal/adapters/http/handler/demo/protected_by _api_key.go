package demohandler

import (
	httpctx "go-archetype/internal/adapters/http/context"

	"github.com/gofiber/fiber/v2"
)

func (h *DemoHandler) ProtectedByAPIKey(c *fiber.Ctx) error {
	log := httpctx.Get(c, h.log)
	log.Info("Hello from protected route by API key!")
	return c.SendString("Hello from protected route by API key!")
}
