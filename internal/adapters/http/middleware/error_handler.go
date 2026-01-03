package middleware

import (
	"errors"

	"go-archetype/internal/adapters/http/context"
	"go-archetype/internal/adapters/http/dto/response"
	"go-archetype/internal/pkg/apperror"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is the global Fiber error handler
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		rid := httpctx.GetRequestID(c)

		// 1️⃣ Application-level error (AppError)
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case apperror.CodeValidation:
				return c.Status(fiber.StatusBadRequest).
					JSON(response.FailMessage(appErr.Message, rid))

			case apperror.CodeNotFound:
				return c.Status(fiber.StatusNotFound).
					JSON(response.FailMessage(appErr.Message, rid))

			case apperror.CodeUnauthorized:
				return c.Status(fiber.StatusUnauthorized).
					JSON(response.FailMessage(appErr.Message, rid))

			case apperror.CodeConflict:
				return c.Status(fiber.StatusConflict).
					JSON(response.FailMessage(appErr.Message, rid))

			default:
				return c.Status(fiber.StatusInternalServerError).
					JSON(response.FailMessage("internal server error", rid))
			}
		}

		// 2️⃣ Fiber native error (from middleware like AuthJWT / AuthAPIKey)
		var fe *fiber.Error
		if errors.As(err, &fe) {
			return c.Status(fe.Code).
				JSON(response.FailMessage(fe.Message, rid))
		}

		// 3️⃣ Unknown / panic / raw error
		return c.Status(fiber.StatusInternalServerError).
			JSON(response.FailMessage("internal server error", rid))
	}
}
