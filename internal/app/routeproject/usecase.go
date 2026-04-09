package routeproject

import (
	"context"
	"math/rand"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
)

type RouteProjectUseCase struct {
	dynamoSvc                dynamoservice.StoredNotificationService
	currentDate              *time.Time
	thankYouJitterMaxMinutes int
}

func NewRouteProjectUseCase(dynamoSvc dynamoservice.StoredNotificationService, currentDate *time.Time, thankYouJitterMaxMinutes int) *RouteProjectUseCase {
	return &RouteProjectUseCase{dynamoSvc: dynamoSvc, currentDate: currentDate, thankYouJitterMaxMinutes: thankYouJitterMaxMinutes}
}

func (u *RouteProjectUseCase) Execute(ctx context.Context, project domain.Project) (domain.ProjectNotification, domain.NotificationType, time.Time, error) {
	existingNotification, err := u.dynamoSvc.GetProjectNotification(ctx, project)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, time.Time{}, err
	}

	now := time.Now().UTC()
	if u.currentDate != nil {
		now = *u.currentDate
	}
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	messageType, targetSendTime, err := computeNotificationType(now, project.Date, project.EndDateTime, project.Status, project.IsTeamLeader, existingNotification, u.thankYouJitterMaxMinutes)
	if err != nil {
		return domain.ProjectNotification{}, domain.Welcome, time.Time{}, err
	}

	if existingNotification == nil {
		return domain.ProjectNotification{
			Name: project.Name,
			Date: project.Date,
			Id:   project.Id,
		}, messageType, targetSendTime, nil
	}
	return *existingNotification, messageType, targetSendTime, nil
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

// AllNotificationsSent is returned when all applicable notifications have already been sent.
// Step Functions catches this by type name to route to EndProjectIteration.
type AllNotificationsSent struct{}

func (e *AllNotificationsSent) Error() string { return "all notifications already sent for project" }

// NotificationsDisabled is returned when ShouldStopNotify is set for the project.
// Step Functions catches this by type name to route to EndProjectIteration.
type NotificationsDisabled struct{}

func (e *NotificationsDisabled) Error() string { return "notifications are disabled for this project" }

func computeNotificationType(now, projectDate, endDateTime time.Time, status string, isTeamLeader bool, existingNotification *domain.ProjectNotification, jitterMaxMinutes int) (domain.NotificationType, time.Time, error) {
	if !isTeamLeader {
		return domain.Welcome, time.Time{}, &NotTeamLeader{}
	}

	if status == "Canceled" {
		return domain.Welcome, time.Time{}, &ProjectCancelled{}
	}

	// Project has not started yet — check welcome/reminder windows
	if now.Before(projectDate) {
		if !shouldSendWelcome(now, projectDate) && !shouldSendReminder(now, projectDate) {
			return domain.Welcome, time.Time{}, &ProjectTooFar{}
		}

		if existingNotification == nil {
			if shouldSendWelcome(now, projectDate) {
				return domain.Welcome, time.Time{}, nil
			}
			if shouldSendReminder(now, projectDate) {
				return domain.Reminder, time.Time{}, nil
			}
		}

		if existingNotification.ShouldStopNotify {
			return domain.Welcome, time.Time{}, &NotificationsDisabled{}
		}

		if !existingNotification.HasSentWelcome {
			if shouldSendWelcome(now, projectDate) {
				return domain.Welcome, time.Time{}, nil
			}
		}

		if !existingNotification.HasSentReminder {
			if shouldSendReminder(now, projectDate) {
				return domain.Reminder, time.Time{}, nil
			}
		}

		return domain.Welcome, time.Time{}, &AllNotificationsSent{}
	}

	// Project is today — schedule thank-you
	if now.Equal(projectDate) {
		if existingNotification == nil || !existingNotification.HasSentThankYou {
			jitter := time.Duration(0)
			if jitterMaxMinutes > 0 {
				jitter = time.Duration(rand.Intn(jitterMaxMinutes+1)) * time.Minute
			}
			targetSendTime := endDateTime.Add(1*time.Hour + jitter)
			return domain.ThankYou, targetSendTime, nil
		}
	}

	return domain.Welcome, time.Time{}, &AllNotificationsSent{}
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
