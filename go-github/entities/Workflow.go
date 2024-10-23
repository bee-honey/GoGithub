package entities

// Struct to parse the workflow runs response
type WorkflowRun struct {
	HeadSha    string `json:"head_sha"`
	Conclusion string `json:"conclusion"`
	Status     string `json:"status"`
}
