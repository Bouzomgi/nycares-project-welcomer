package models

type RecordMessageInput = RequestApprovalInput

type RecordMessageOutput struct {
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageToSend               message             `json:"message"`
	RecordedProjectNotification projectNotification `json:"recordedProjectNotification"`
	ExecutionId                 string              `json:"executionId"`
}
