package response

type Error struct {
	Message   string      `json:"message"`
	Error     interface{} `json:"error,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}
