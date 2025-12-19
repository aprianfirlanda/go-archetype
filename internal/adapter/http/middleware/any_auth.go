package middleware

import "github.com/gofiber/fiber/v2"

// AnyAuth allows request if ANY middleware succeeds.
func AnyAuth(middlewares ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		for _, mw := range middlewares {
			// clone ctx to avoid side effects
			ctx := c.Context()
			if err := mw(c); err == nil {
				// auth passed
				return c.Next()
			}
			// reset ctx if needed (Fiber v2 safe here)
			c.SetUserContext(ctx)
		}
		return fiber.ErrUnauthorized
	}
}
