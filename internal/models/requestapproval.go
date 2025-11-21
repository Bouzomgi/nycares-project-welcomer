package models

type RequestApprovalInput struct {
	TaskToken     string  `json:"taskToken"`
	Auth          Auth    `json:"auth"`
	Project       project `json:"project"`
	MessageToSend message `json:"message"`
}

type RequestApprovalOutput struct {
	Auth          Auth    `json:"auth"`
	Project       project `json:"project"`
	MessageToSend message `json:"message"`
}
