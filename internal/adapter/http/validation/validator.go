package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			var msgs []string
			for _, fe := range ve {
				msgs = append(msgs, fmt.Sprintf(
					"%s failed on '%s'",
					strings.ToLower(fe.Field()),
					fe.Tag(),
				))
			}
			return fmt.Errorf(strings.Join(msgs, ", "))
		}
		return err
	}
	return nil
}
