package computemessage

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
)

type ComputeMessageUseCase struct {
	dynamoSvc dynamoservice.StoredNotificationService
}

func NewComputeMessageUseCase(dynamoSvc dynamoservice.StoredNotificationService) *ComputeMessageUseCase {
	return &ComputeMessageUseCase{dynamoSvc}
}

func (u *ComputeMessageUseCase) Execute(ctx context.Context, messageBucketName string, project domain.Project) (domain.ProjectNotification, domain.NotificationType, string, error) {

	existingNotification, err := u.dynamoSvc.GetProjectNotification(ctx, project)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	messageType, err := computeNotificationType(project.Date, existingNotification)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	s3MessageRef, err := computeS3MessageRefPath(messageBucketName, project.Name, messageType)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	return *existingNotification, messageType, s3MessageRef, nil
}

func computeNotificationType(projectDate time.Time, existingNotification *domain.ProjectNotification) (domain.NotificationType, error) {
	if existingNotification == nil {
		if shouldSendWelcome(time.Now(), projectDate) {
			return domain.Welcome, nil
		}
	}

	if existingNotification.ShouldStopNotify {
		return domain.Welcome, fmt.Errorf("notifications disabled for project")
	}

	if !existingNotification.HasSentWelcome {
		if shouldSendWelcome(time.Now(), projectDate) {
			return domain.Welcome, nil
		}
	}

	if !existingNotification.HasSentReminder {
		if shouldSendReminder(time.Now(), projectDate) {
			return domain.Reminder, nil
		}
	}

	return domain.Welcome, fmt.Errorf("all notifications already sent")
}

// ShouldSendWelcome checks if a welcome message should be sent
func shouldSendWelcome(now, projectDate time.Time) bool {
	weekBefore := projectDate.AddDate(0, 0, -7)
	return now.After(weekBefore) && now.Before(projectDate)
}

// ShouldSendReminder checks if a reminder message should be sent
func shouldSendReminder(now, projectDate time.Time) bool {
	weekBefore := projectDate.AddDate(0, 0, -2)
	return now.After(weekBefore) && now.Before(projectDate)
}

func computeS3MessageRefPath(s3BucketName, projectName string, messageType domain.NotificationType) (string, error) {
	messageTypeStr := strings.ToLower(messageType.String())
	s3MessageRef := fmt.Sprintf("s3://%s/%s/%s.md", s3BucketName, toKebabCase(projectName), messageTypeStr)

	if isValidS3URI(s3MessageRef) {
		return s3MessageRef, nil
	}

	return "", fmt.Errorf("could not compute valid s3 bucket reference for message")
}

func toKebabCase(s string) string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}

	for i := range words {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "-")
}

var basicS3URIRegex = regexp.MustCompile(`^s3://[a-z0-9\-]{3,63}/.+$`)

func isValidS3URI(uri string) bool {
	return basicS3URIRegex.MatchString(uri)
}
