package models

type RouteProjectInput struct {
	Auth        Auth    `json:"auth"`
	Project     project `json:"project"`
	ExecutionId string  `json:"executionId"`
}

type RouteProjectOutput struct {
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageType                 string              `json:"messageType"`
	ExecutionId                 string              `json:"executionId"`
}
