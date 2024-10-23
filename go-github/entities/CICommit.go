package entities

// Struct to parse the commit response
type CICommit struct {
	Sha        string `json:"sha"`
	Message    string `json:"message"`
	Author     string `json:"author"`
	Status     string `json:"status"` // Includes the status of the associated workflow run
	Conclusion string `json:"conclusion"`
}
