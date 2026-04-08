package dto

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	ProjectId        string `json:"projectId"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	HasSentThankYou  bool   `json:"hasSentThankYou"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
	LastUpdated      string `json:"lastUpdated"`
}
