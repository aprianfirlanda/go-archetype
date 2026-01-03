package apperror

func Validation(message string, err error) *AppError {
	return New(CodeValidation, message, err)
}

func NotFound(message string, err error) *AppError {
	return New(CodeNotFound, message, err)
}

func Conflict(message string, err error) *AppError {
	return New(CodeConflict, message, err)
}

func Unauthorized(message string, err error) *AppError {
	return New(CodeUnauthorized, message, err)
}

func Internal(message string, err error) *AppError {
	return New(CodeInternal, message, err)
}
