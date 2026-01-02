package taskresp

type BulkDeleteValidateError struct {
	IDs []string `json:"ids,omitempty"`
}

type BulkDelete struct {
	Deleted []string            `json:"deleted"`
	Failed  []BulkDeleteFailure `json:"failed"`
}

type BulkDeleteFailure struct {
	PublicID string `json:"id"`
	Reason   string `json:"reason"`
}
