package taskresult

type BulkDeleteResult struct {
	Deleted []string            `json:"deleted"`
	Failed  []BulkDeleteFailure `json:"failed"`
}

type BulkDeleteFailure struct {
	PublicID string `json:"id"`
	Reason   string `json:"reason"`
}
