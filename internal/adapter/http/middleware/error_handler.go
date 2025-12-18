package middleware

import (
	"errors"
	"go-archetype/internal/adapter/http/context"
	"go-archetype/internal/adapter/http/dto/response"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler returns a Fiber ErrorHandler that converts errors to JSON responses.
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Default to 500
		status := fiber.StatusInternalServerError

		// If it's a *fiber.Error, use its status code
		var fe *fiber.Error
		if errors.As(err, &fe) {
			status = fe.Code
		}

		// Get request ID (from requestid middleware)
		rid := context.GetRequestID(c)

		// Build JSON error response
		resp := response.ErrorResponse{
			Message:   err.Error(),
			RequestID: rid,
		}

		// Ensure the content type is JSON and send a response
		c.Type("json", "utf-8")
		return c.Status(status).JSON(resp)
	}
}
