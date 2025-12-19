package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FieldErrors map[string][]string

func ValidateStruct(s any) (FieldErrors, error) {
	err := validate.Struct(s)
	if err == nil {
		return nil, nil
	}

	var ve validator.ValidationErrors
	ok := errors.As(err, &ve)
	if !ok {
		return nil, err
	}

	errs := FieldErrors{}

	for _, fe := range ve {
		field := toSnakeCase(fe.Field())

		errs[field] = append(errs[field], messageFor(fe))
	}

	return errs, nil
}
