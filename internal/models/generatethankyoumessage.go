package models

// GenerateThankYouMessageInput comes directly from RouteProjectOutput.
type GenerateThankYouMessageInput = RouteProjectOutput

// GenerateThankYouMessageOutput matches ComputeMessageOutput so that
// ScheduleThankYou can read $.message.type downstream.
type GenerateThankYouMessageOutput = ComputeMessageOutput
