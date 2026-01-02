package taskreq

type UpdateStatus struct {
	Status string `json:"status" validate:"required,oneof=todo in_progress done"`
}
