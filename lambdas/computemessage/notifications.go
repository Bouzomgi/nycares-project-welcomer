package main

import (
	"fmt"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	log "github.com/sirupsen/logrus"
)

func HandleProjectNotification(s3BucketName string, project models.Project, existingNotification *models.ProjectNotification) (*models.ProjectNotification, models.SendableMessage, error) {
	if s3BucketName == "" {
		log.Errorf("S3 bucket name is empty for project %s", project.Name)
		return nil, models.SendableMessage{}, fmt.Errorf("S3 bucket name is required")
	}

	if existingNotification == nil {
		log.Infof("No existing notification for project %s, creating welcome message", project.Name)
		notification := createNewNotification(s3BucketName, project, models.Welcome)
		return nil, notification, nil
	}

	if existingNotification.ShouldStopNotify {
		log.Infof("Notifications are disabled for project %s", project.Name)
		return nil, models.SendableMessage{}, fmt.Errorf("ShouldStopNotify is true")
	}

	if !existingNotification.HasSentWelcome {
		log.Infof("Welcome not sent yet for project %s, creating welcome message", project.Name)
		notification := createNewNotification(s3BucketName, project, models.Welcome)
		return existingNotification, notification, nil
	}

	if !existingNotification.HasSentReminder {
		log.Infof("Reminder not sent yet for project %s, creating reminder message", project.Name)
		notification := createNewNotification(s3BucketName, project, models.Reminder)
		return existingNotification, notification, nil
	}

	log.Infof("All notifications already sent for project %s", project.Name)
	return nil, models.SendableMessage{}, fmt.Errorf("AllNotificationsAlreadySent")
}

func createNewNotification(s3BucketName string, project models.Project, messageType models.MessageType) models.SendableMessage {
	messageTypeStr := strings.ToLower(messageType.String())
	s3Destination := fmt.Sprintf("s3://%s/%s/%s.md", s3BucketName, ToCamelCase(project.Name), messageTypeStr)

	log.Infof("Created notification for project %s: %s", project.Name, s3Destination)

	return models.SendableMessage{
		Type:        messageTypeStr,
		TemplateRef: s3Destination,
	}
}
