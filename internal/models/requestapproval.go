package models

type RequestApprovalInput struct {
	TaskToken                   string              `json:"taskToken"`
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageToSend               message             `json:"message"`
	ExecutionId                 string              `json:"executionId"`
}

type RequestApprovalOutput struct {
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageToSend               message             `json:"message"`
	ExecutionId                 string              `json:"executionId"`
}
