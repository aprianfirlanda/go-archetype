package apperror

// Code represents application-level error codes
type Code string

const (
	CodeNotFound     Code = "NOT_FOUND"
	CodeValidation   Code = "VALIDATION_ERROR"
	CodeConflict     Code = "CONFLICT"
	CodeUnauthorized Code = "UNAUTHORIZED"
	CodeInternal     Code = "INTERNAL_ERROR"
)

// AppError is the unified application error
type AppError struct {
	Code    Code
	Message string
	Err     error // wrapped original error (optional)
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// Unwrap returns the wrapped error. it automatically unwraps nested errors on errors.Is
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError
func New(code Code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
