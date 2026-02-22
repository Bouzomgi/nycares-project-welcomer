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
	dynamoSvc   dynamoservice.StoredNotificationService
	currentDate *time.Time
}

func NewComputeMessageUseCase(dynamoSvc dynamoservice.StoredNotificationService, currentDate *time.Time) *ComputeMessageUseCase {
	return &ComputeMessageUseCase{dynamoSvc: dynamoSvc, currentDate: currentDate}
}

func (u *ComputeMessageUseCase) Execute(ctx context.Context, messageBucketName string, project domain.Project) (domain.ProjectNotification, domain.NotificationType, string, error) {

	existingNotification, err := u.dynamoSvc.GetProjectNotification(ctx, project)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	now := time.Now().UTC()
	if u.currentDate != nil {
		now = *u.currentDate
	}
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	messageType, err := computeNotificationType(now, project.Date, existingNotification)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	s3MessageRef, err := computeS3MessageRefPath(messageBucketName, project.Name, messageType)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, "", err
	}

	if existingNotification == nil {
		return domain.ProjectNotification{
			Name: project.Name,
			Date: project.Date,
			Id:   project.Id,
		}, messageType, s3MessageRef, nil
	}
	return *existingNotification, messageType, s3MessageRef, nil
}

// ProjectTooFar is returned when the project is too far in the future for any notification.
// Step Functions catches this by type name to route to EndProjectIteration.
type ProjectTooFar struct{}

func (e *ProjectTooFar) Error() string { return "project is too far in the future to notify" }

// ProjectPassed is returned when the project date has already passed.
// Step Functions catches this by type name to route to EndProjectIteration.
type ProjectPassed struct{}

func (e *ProjectPassed) Error() string { return "project date has already passed" }

func computeNotificationType(now, projectDate time.Time, existingNotification *domain.ProjectNotification) (domain.NotificationType, error) {
	if !now.Before(projectDate) {
		return domain.Welcome, &ProjectPassed{}
	}

	if !shouldSendWelcome(now, projectDate) && !shouldSendReminder(now, projectDate) {
		return domain.Welcome, &ProjectTooFar{}
	}

	if existingNotification == nil {
		if shouldSendWelcome(now, projectDate) {
			return domain.Welcome, nil
		}
		if shouldSendReminder(now, projectDate) {
			return domain.Reminder, nil
		}
	}

	if existingNotification.ShouldStopNotify {
		return domain.Welcome, fmt.Errorf("notifications disabled for project")
	}

	if !existingNotification.HasSentWelcome {
		if shouldSendWelcome(now, projectDate) {
			return domain.Welcome, nil
		}
	}

	if !existingNotification.HasSentReminder {
		if shouldSendReminder(now, projectDate) {
			return domain.Reminder, nil
		}
	}

	return domain.Welcome, fmt.Errorf("all notifications already sent")
}

const (
	welcomeLeadDays  = 7
	reminderLeadDays = 2
)

func shouldSendWelcome(now, projectDate time.Time) bool {
	cutoff := projectDate.AddDate(0, 0, -welcomeLeadDays)
	return now.After(cutoff) && now.Before(projectDate)
}

func shouldSendReminder(now, projectDate time.Time) bool {
	cutoff := projectDate.AddDate(0, 0, -reminderLeadDays)
	return now.After(cutoff) && now.Before(projectDate)
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
