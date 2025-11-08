package models

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type Auth struct {
	Cookies []Cookie `json:"cookies"`
}

type Project struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
	LastUpdated      string `json:"lastUpdated"`
}

type SendableMessage struct {
	Type        string `json:"type"`
	TemplateRef string `json:"templateRef"`
}

type MessageType int

const (
	Welcome MessageType = iota
	Reminder
)

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
