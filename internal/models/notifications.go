package models

type SendableMessage struct {
	Type        string `json:"type"`
	TemplateRef string `json:"templateRef"`
}

type MessageType int

const (
	Welcome MessageType = iota
	Reminder
)

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
	LastUpdated      string `json:"lastUpdated"`
}

func (m MessageType) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case Reminder:
		return "reminder"
	default:
		return "unknown"
	}
}
