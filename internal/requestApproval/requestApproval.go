package requestapproval

import (
	"fmt"
)

func CreateNotificationMessage(messageType, projectName, projectDate, messageTemplateRef string) string {
	return fmt.Sprintf(
		"Requesting approval to publish %s message to project %s on %s with template %s",
		messageType,
		projectName,
		projectDate,
		messageTemplateRef,
	)
}

func CreateNotificationSubject(messageType, projectName, projectDate string) string {
	return fmt.Sprintf(
		"Send %s to %s on %s?",
		messageType,
		projectName,
		projectDate,
	)
}
