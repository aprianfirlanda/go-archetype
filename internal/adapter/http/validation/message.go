package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("minimum value is %s", fe.Param())
	case "max":
		return fmt.Sprintf("maximum value is %s", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of [%s]", fe.Param())
	case "email":
		return "must be a valid email"
	default:
		return "is invalid"
	}
}
