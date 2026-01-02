package taskresp

type UpdateValidateError struct {
	Title       []string `json:"title,omitempty"`
	Description []string `json:"description,omitempty"`
	Priority    []string `json:"priority,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
