package computemessage

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

func TestToKebabCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Park Cleanup", "park-cleanup"},
		{"Food Bank Volunteer", "food-bank-volunteer"},
		{"UPPER CASE", "upper-case"},
		{"single", "single"},
		{"", ""},
		{"  extra   spaces  ", "extra-spaces"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toKebabCase(tt.input)
			if got != tt.expected {
				t.Errorf("toKebabCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestIsValidS3URI(t *testing.T) {
	tests := []struct {
		uri      string
		expected bool
	}{
		{"s3://my-bucket/path/to/file.md", true},
		{"s3://bucket123/file", true},
		{"s3://ab/file", false}, // bucket name too short
		{"http://bucket/file", false},
		{"s3:///file", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			got := isValidS3URI(tt.uri)
			if got != tt.expected {
				t.Errorf("isValidS3URI(%q) = %v, want %v", tt.uri, got, tt.expected)
			}
		})
	}
}

func TestComputeS3MessageRefPath(t *testing.T) {
	ref, err := computeS3MessageRefPath("my-bucket", "Park Cleanup", domain.Welcome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "s3://my-bucket/park-cleanup/welcome.md"
	if ref != expected {
		t.Errorf("got %q, want %q", ref, expected)
	}

	ref, err = computeS3MessageRefPath("my-bucket", "Food Bank", domain.Reminder)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected = "s3://my-bucket/food-bank/reminder.md"
	if ref != expected {
		t.Errorf("got %q, want %q", ref, expected)
	}
}

func TestComputeNotificationType(t *testing.T) {
	projectDate := time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC)

	notif := func(welcome, reminder, stop bool) *domain.ProjectNotification {
		return &domain.ProjectNotification{
			HasSentWelcome:   welcome,
			HasSentReminder:  reminder,
			ShouldStopNotify: stop,
		}
	}

	tests := []struct {
		name        string
		now         time.Time
		status      string
		existing    *domain.ProjectNotification
		wantType    domain.NotificationType
		wantErrType error
	}{
		{
			name:        "project is cancelled",
			now:         projectDate.AddDate(0, 0, -5),
			status:      "Canceled",
			wantErrType: &ProjectCancelled{},
		},
		{
			name:        "project date has passed",
			now:         projectDate,
			wantErrType: &ProjectPassed{},
		},
		{
			name:        "project is in the past",
			now:         projectDate.AddDate(0, 0, 1),
			wantErrType: &ProjectPassed{},
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
			existing:    notif(true, false, true),
			wantErrType: &NotificationsDisabled{},
		},
		{
			name:     "welcome not sent, in welcome window",
			now:      projectDate.AddDate(0, 0, -5),
			existing: notif(false, false, false),
			wantType: domain.Welcome,
		},
		{
			name:     "welcome sent, reminder not sent, in reminder window",
			now:      projectDate.AddDate(0, 0, -1),
			existing: notif(true, false, false),
			wantType: domain.Reminder,
		},
		{
			name:        "welcome sent, not in reminder window yet",
			now:         projectDate.AddDate(0, 0, -5),
			existing:    notif(true, false, false),
			wantErrType: &AllNotificationsSent{},
		},
		{
			name:        "both welcome and reminder already sent",
			now:         projectDate.AddDate(0, 0, -1),
			existing:    notif(true, true, false),
			wantErrType: &AllNotificationsSent{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, err := computeNotificationType(tt.now, projectDate, tt.status, true, tt.existing)

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
		})
	}
}
