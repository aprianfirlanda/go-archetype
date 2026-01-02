package taskreq

type BulkDelete struct {
	IDs []string `json:"ids" validate:"required,min=1,dive,required"`
}
