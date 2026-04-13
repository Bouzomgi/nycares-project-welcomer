package models

// GenerateThankYouMessageInput extends RouteProjectOutput with an optional
// RefinementContext supplied when the approver requests a refined regeneration.
type GenerateThankYouMessageInput struct {
	Auth                        Auth                `json:"auth"`
	ExistingProjectNotification projectNotification `json:"existingProjectNotification"`
	MessageType                 string              `json:"messageType"`
	TargetSendTime              string              `json:"targetSendTime,omitempty"`
	ExecutionId                 string              `json:"executionId"`
	RefinementContext           string              `json:"refinementContext,omitempty"`
}

// GenerateThankYouMessageOutput matches ComputeMessageOutput so that
// ScheduleThankYou can read $.message.type downstream.
type GenerateThankYouMessageOutput = ComputeMessageOutput
