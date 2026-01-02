package taskresp

type UpdateStatusValidateError struct {
	Status []string `json:"status,omitempty"`
}
