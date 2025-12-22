package response

type ErrorResponse struct {
	Message   string      `json:"message"`
	RequestID string      `json:"requestId,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
}
