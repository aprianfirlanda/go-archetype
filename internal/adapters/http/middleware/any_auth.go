package middleware

import (
	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/response"

	"github.com/gofiber/fiber/v2"
)

// AnyAuth allows request if ANY middleware succeeds
func AnyAuth(middlewares ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, mw := range middlewares {
			if err := mw(c); err == nil {
				return c.Next()
			}
		}

		// All middlewares failed
		rid := httpctx.GetRequestID(c)
		resp := response.FailMessage("invalid or missing authentication credentials", rid)

		c.Type("json", "utf-8")
		return c.Status(fiber.StatusUnauthorized).JSON(resp)
	}
}
