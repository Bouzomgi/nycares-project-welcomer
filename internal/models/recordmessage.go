package models

type RecordMessageInput = RequestApprovalInput

type RecordMessageOutput struct {
	TaskToken                   string              `json:"taskToken"`
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageToSend               message             `json:"message"`
	RecordedProjectNotification projectNotification `json:"recordedProjectNotification"`
}
