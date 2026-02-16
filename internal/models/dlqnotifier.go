package models

// DLQNotifierInput represents the error payload from Step Functions catch blocks.
// Step Functions always includes Error and Cause fields.
type DLQNotifierInput struct {
	Error string `json:"Error"`
	Cause string `json:"Cause"`
}
