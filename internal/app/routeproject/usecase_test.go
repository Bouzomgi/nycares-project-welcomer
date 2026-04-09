package routeproject

import (
	"fmt"
	"testing"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

func TestShouldSendWelcome(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{"8 days before - too early", projectDate.AddDate(0, 0, -8), false},
		{"7 days before - boundary", projectDate.AddDate(0, 0, -7).Add(time.Hour), true},
		{"5 days before - in window", projectDate.AddDate(0, 0, -5), true},
		{"1 day before - in window", projectDate.AddDate(0, 0, -1), true},
		{"project day - too late", projectDate, false},
		{"after project - too late", projectDate.AddDate(0, 0, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSendWelcome(tt.now, projectDate)
			if got != tt.expected {
				t.Errorf("shouldSendWelcome(%v, %v) = %v, want %v", tt.now, projectDate, got, tt.expected)
			}
		})
	}
}

func TestShouldSendReminder(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		now      time.Time
		expected bool
	}{
		{"3 days before - too early", projectDate.AddDate(0, 0, -3), false},
		{"2 days before - boundary", projectDate.AddDate(0, 0, -2).Add(time.Hour), true},
		{"1 day before - in window", projectDate.AddDate(0, 0, -1), true},
		{"project day - too late", projectDate, false},
		{"after project - too late", projectDate.AddDate(0, 0, 1), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSendReminder(tt.now, projectDate)
			if got != tt.expected {
				t.Errorf("shouldSendReminder(%v, %v) = %v, want %v", tt.now, projectDate, got, tt.expected)
			}
		})
	}
}

func TestComputeNotificationType(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)
	endDateTime := time.Date(2025, 3, 15, 14, 0, 0, 0, time.UTC)

	notif := func(welcome, reminder, thankYou, stop bool) *domain.ProjectNotification {
		return &domain.ProjectNotification{
			HasSentWelcome:   welcome,
			HasSentReminder:  reminder,
			HasSentThankYou:  thankYou,
			ShouldStopNotify: stop,
		}
	}

	tests := []struct {
		name              string
		now               time.Time
		status            string
		existing          *domain.ProjectNotification
		wantType          domain.NotificationType
		wantErrType       error
		wantTargetNonZero bool
	}{
		{
			name:        "project is cancelled",
			now:         projectDate.AddDate(0, 0, -5),
			status:      "Canceled",
			wantErrType: &ProjectCancelled{},
		},
		{
			name:        "project too far in future",
			now:         projectDate.AddDate(0, 0, -30),
			wantErrType: &ProjectTooFar{},
		},
		{
			name:     "no existing notification, in welcome window",
			now:      projectDate.AddDate(0, 0, -5),
			existing: nil,
			wantType: domain.Welcome,
		},
		{
			name:     "no existing notification, in reminder window",
			now:      projectDate.AddDate(0, 0, -1),
			existing: nil,
			wantType: domain.Welcome, // welcome takes priority for first notification
		},
		{
			name:        "notifications disabled",
			now:         projectDate.AddDate(0, 0, -5),
			existing:    notif(true, false, false, true),
			wantErrType: &NotificationsDisabled{},
		},
		{
			name:     "welcome not sent, in welcome window",
			now:      projectDate.AddDate(0, 0, -5),
			existing: notif(false, false, false, false),
			wantType: domain.Welcome,
		},
		{
			name:     "welcome sent, reminder not sent, in reminder window",
			now:      projectDate.AddDate(0, 0, -1),
			existing: notif(true, false, false, false),
			wantType: domain.Reminder,
		},
		{
			name:        "welcome sent, not in reminder window yet",
			now:         projectDate.AddDate(0, 0, -5),
			existing:    notif(true, false, false, false),
			wantErrType: &AllNotificationsSent{},
		},
		{
			name:        "both welcome and reminder already sent, before project",
			now:         projectDate.AddDate(0, 0, -1),
			existing:    notif(true, true, false, false),
			wantErrType: &AllNotificationsSent{},
		},
		{
			name:              "project day, no thank-you sent yet",
			now:               projectDate,
			existing:          nil,
			wantType:          domain.ThankYou,
			wantTargetNonZero: true,
		},
		{
			name:        "project in past, no thank-you sent",
			now:         projectDate.AddDate(0, 0, 1),
			existing:    notif(true, true, false, false),
			wantErrType: &AllNotificationsSent{},
		},
		{
			name:        "project in past, thank-you already sent",
			now:         projectDate.AddDate(0, 0, 1),
			existing:    notif(true, true, true, false),
			wantErrType: &AllNotificationsSent{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotTarget, err := computeNotificationType(tt.now, projectDate, endDateTime, tt.status, true, tt.existing)

			if tt.wantErrType != nil {
				if err == nil {
					t.Fatalf("expected error %T, got nil", tt.wantErrType)
				}
				if got, want := err, tt.wantErrType; fmt.Sprintf("%T", got) != fmt.Sprintf("%T", want) {
					t.Errorf("error type = %T, want %T", got, want)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotType != tt.wantType {
				t.Errorf("notification type = %v, want %v", gotType, tt.wantType)
			}
			if tt.wantTargetNonZero && gotTarget.IsZero() {
				t.Errorf("expected non-zero targetSendTime, got zero")
			}
			if !tt.wantTargetNonZero && !gotTarget.IsZero() {
				t.Errorf("expected zero targetSendTime, got %v", gotTarget)
			}
		})
	}
}
