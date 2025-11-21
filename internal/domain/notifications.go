package domain

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
}

type NotificationType int

const (
	Welcome NotificationType = iota
	Reminder
)

func (m NotificationType) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case Reminder:
		return "reminder"
	default:
		return "unknown"
	}
}
