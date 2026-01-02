package taskresult

type BulkUpdateStatusResult struct {
	Updated []string               `json:"updated"`
	Failed  []BulkUpdateStatusFail `json:"failed"`
}

type BulkUpdateStatusFail struct {
	PublicID string `json:"id"`
	Reason   string `json:"reason"`
}
