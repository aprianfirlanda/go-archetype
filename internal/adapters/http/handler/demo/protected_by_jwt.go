package demohandler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *DemoHandler) ProtectedByJWT(c *fiber.Ctx) error {
	return c.SendString("Hello from protected route by JWT!")
}
