package taskreq

type BulkUpdateStatus struct {
	IDs    []string `json:"ids" validate:"required,min=1,dive,required"`
	Status string   `json:"status" validate:"required,oneof=todo in_progress done"`
}
