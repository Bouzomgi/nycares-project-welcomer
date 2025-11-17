package computemessage

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/s3service"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// ComputeProjectMessage determines what message should be sent for a project
func ComputeProjectMessage(s3client *s3service.S3Service, project models.Project, existingNotification *models.ProjectNotification) (models.SendableMessage, error) {
	if existingNotification == nil {
		notification := createNewNotification(s3client, project, models.Welcome)
		return notification, nil
	}

	if existingNotification.ShouldStopNotify {
		return models.SendableMessage{}, fmt.Errorf("notifications disabled for project")
	}

	if !existingNotification.HasSentWelcome {
		notification := createNewNotification(s3client, project, models.Welcome)
		return notification, nil
	}

	if !existingNotification.HasSentReminder {
		notification := createNewNotification(s3client, project, models.Reminder)
		return notification, nil
	}

	return models.SendableMessage{}, fmt.Errorf("all notifications already sent")
}

func createNewNotification(s3Service *s3service.S3Service, project models.Project, messageType models.MessageType) models.SendableMessage {
	messageTypeStr := strings.ToLower(messageType.String())
	projectRefSuffix := fmt.Sprintf("%s/%s.md", toCamelCase(project.Name), messageTypeStr)
	s3Destination := s3Service.CreateS3Path(projectRefSuffix)

	return models.SendableMessage{
		Type:        messageTypeStr,
		TemplateRef: s3Destination,
	}
}

func toCamelCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	caser := cases.Title(language.English)

	result := strings.ToLower(words[0])
	for i := 1; i < len(words); i++ {
		w := strings.ToLower(words[i])
		result += caser.String(w)
	}

	return result
}

// ShouldSendWelcome checks if a welcome message should be sent
func ShouldSendWelcome(projectDate time.Time) bool {
	now := time.Now()
	weekBefore := projectDate.AddDate(0, 0, -7)
	return now.After(weekBefore) && now.Before(projectDate)
}

// ShouldSendReminder checks if a reminder message should be sent
func ShouldSendReminder(projectDate time.Time) bool {
	now := time.Now()
	twoDaysBefore := projectDate.AddDate(0, 0, -2)
	return now.After(twoDaysBefore) && now.Before(projectDate)
}
