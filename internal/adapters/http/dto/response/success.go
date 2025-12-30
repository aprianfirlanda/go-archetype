package response

type Success struct {
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Meta      interface{} `json:"meta,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}
