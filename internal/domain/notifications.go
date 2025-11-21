package domain

import "fmt"

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

func ParseNotificationType(s string) (NotificationType, error) {
	switch s {
	case "welcome":
		return Welcome, nil
	case "reminder":
		return Reminder, nil
	default:
		return 0, fmt.Errorf("unknown notification type: %s", s)
	}
}
