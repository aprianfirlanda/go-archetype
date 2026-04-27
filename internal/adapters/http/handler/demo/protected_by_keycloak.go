package demohandler

import "github.com/gofiber/fiber/v2"

func (h *DemoHandler) ProtectedByKeycloak(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Hello from protected route by Keycloak!",
	})
}
