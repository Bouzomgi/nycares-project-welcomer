package domain

import (
	"fmt"
	"time"
)

// Maps to Dynamo Schema
type ProjectNotification struct {
	Name             string    `json:"name"`
	Date             time.Time `json:"date"`
	Id               string    `json:"id"`
	HasSentWelcome   bool      `json:"hasSentWelcome"`
	HasSentReminder  bool      `json:"hasSentReminder"`
	HasSentThankYou  bool      `json:"hasSentThankYou"`
	ShouldStopNotify bool      `json:"shouldStopNotify"`
}

type NotificationType int

const (
	Welcome NotificationType = iota
	Reminder
	ThankYou
)

func (m NotificationType) String() string {
	switch m {
	case Welcome:
		return "welcome"
	case Reminder:
		return "reminder"
	case ThankYou:
		return "thankYou"
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
	case "thankYou":
		return ThankYou, nil
	default:
		return 0, fmt.Errorf("unknown notification type: %s", s)
	}
}
