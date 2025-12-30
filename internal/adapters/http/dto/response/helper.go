package response

func OK(data interface{}, requestID string) Success {
	return Success{
		Data:      data,
		RequestID: requestID,
	}
}

func OKPaginate(data interface{}, meta PaginationMeta, requestID string) Success {
	return Success{
		Data:      data,
		Meta:      meta,
		RequestID: requestID,
	}
}

func OKMessage(message string, requestID string) Success {
	return Success{
		Message:   message,
		RequestID: requestID,
	}
}

func Fail(message string, err interface{}, requestID string) Error {
	return Error{
		Message:   message,
		Error:     err,
		RequestID: requestID,
	}
}

func FailMessage(message, requestID string) Error {
	return Error{
		Message:   message,
		RequestID: requestID,
	}
}
