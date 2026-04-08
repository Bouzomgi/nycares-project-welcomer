package routeproject

import (
	"context"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
)

type RouteProjectUseCase struct {
	dynamoSvc   dynamoservice.StoredNotificationService
	currentDate *time.Time
}

func NewRouteProjectUseCase(dynamoSvc dynamoservice.StoredNotificationService, currentDate *time.Time) *RouteProjectUseCase {
	return &RouteProjectUseCase{dynamoSvc: dynamoSvc, currentDate: currentDate}
}

func (u *RouteProjectUseCase) Execute(ctx context.Context, project domain.Project) (domain.ProjectNotification, domain.NotificationType, error) {
	existingNotification, err := u.dynamoSvc.GetProjectNotification(ctx, project)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, err
	}

	now := time.Now().UTC()
	if u.currentDate != nil {
		now = *u.currentDate
	}
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	messageType, err := computeNotificationType(now, project.Date, project.Status, project.IsTeamLeader, existingNotification)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, err
	}

	if existingNotification == nil {
		return domain.ProjectNotification{
			Name: project.Name,
			Date: project.Date,
			Id:   project.Id,
		}, messageType, nil
	}
	return *existingNotification, messageType, nil
}

// NotTeamLeader is returned when the authenticated user is not the team leader for the project.
// Step Functions catches this by type name to route to EndProjectIteration.
type NotTeamLeader struct{}

func (e *NotTeamLeader) Error() string { return "user is not the team leader for this project" }

// ProjectCancelled is returned when the project status is Cancelled.
// Step Functions catches this by type name to route to EndProjectIteration.
type ProjectCancelled struct{}

func (e *ProjectCancelled) Error() string { return "project is cancelled" }

// ProjectTooFar is returned when the project is too far in the future for any notification.
// Step Functions catches this by type name to route to EndProjectIteration.
type ProjectTooFar struct{}

func (e *ProjectTooFar) Error() string { return "project is too far in the future to notify" }

// ProjectPassed is returned when the project date has already passed.
// Step Functions catches this by type name to route to EndProjectIteration.
type ProjectPassed struct{}

func (e *ProjectPassed) Error() string { return "project date has already passed" }

// AllNotificationsSent is returned when all applicable notifications have already been sent.
// Step Functions catches this by type name to route to EndProjectIteration.
type AllNotificationsSent struct{}

func (e *AllNotificationsSent) Error() string { return "all notifications already sent for project" }

// NotificationsDisabled is returned when ShouldStopNotify is set for the project.
// Step Functions catches this by type name to route to EndProjectIteration.
type NotificationsDisabled struct{}

func (e *NotificationsDisabled) Error() string { return "notifications are disabled for this project" }

func computeNotificationType(now, projectDate time.Time, status string, isTeamLeader bool, existingNotification *domain.ProjectNotification) (domain.NotificationType, error) {
	if !isTeamLeader {
		return domain.Welcome, &NotTeamLeader{}
	}

	if status == "Canceled" {
		return domain.Welcome, &ProjectCancelled{}
	}

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
		return domain.Welcome, &NotificationsDisabled{}
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

	return domain.Welcome, &AllNotificationsSent{}
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
