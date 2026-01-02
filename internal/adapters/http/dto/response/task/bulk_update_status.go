package taskresp

type BulkUpdateStatusValidateError struct {
	IDs    []string `json:"ids,omitempty"`
	Status []string `json:"status,omitempty"`
}

type BulkUpdateStatus struct {
	Updated []string               `json:"updated"`
	Failed  []BulkUpdateStatusFail `json:"failed"`
}

type BulkUpdateStatusFail struct {
	PublicID string `json:"id"`
	Reason   string `json:"reason"`
}
