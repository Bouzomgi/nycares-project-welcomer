package main

import (
	"strings"
	"testing"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

func TestHandleProjectNotification_NoExistingNotification(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
	s3Bucket := "test-bucket"

	row, msg, err := HandleProjectNotification(s3Bucket, project, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row != nil {
		t.Errorf("expected row to be nil, got %+v", row)
	}
	if msg.Type != "welcome" {
		t.Errorf("expected message type 'welcome', got %s", msg.Type)
	}
}

func TestHandleProjectNotification_ShouldStopNotify(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
	s3Bucket := "test-bucket"
	existing := &models.ProjectNotification{ShouldStopNotify: true}

	_, _, err := HandleProjectNotification(s3Bucket, project, existing)
	if err == nil || !strings.Contains(err.Error(), "ShouldStopNotify is true") {
		t.Errorf("expected 'ShouldStopNotify is true' error, got %v", err)
	}
}

func TestHandleProjectNotification_WelcomeNotSent(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
	s3Bucket := "test-bucket"
	existing := &models.ProjectNotification{HasSentWelcome: false}

	row, msg, err := HandleProjectNotification(s3Bucket, project, existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row != existing {
		t.Errorf("expected row to be returned, got %+v", row)
	}
	if msg.Type != "welcome" {
		t.Errorf("expected message type 'welcome', got %s", msg.Type)
	}
}

func TestHandleProjectNotification_ReminderNotSent(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
	s3Bucket := "test-bucket"
	existing := &models.ProjectNotification{HasSentWelcome: true, HasSentReminder: false}

	row, msg, err := HandleProjectNotification(s3Bucket, project, existing)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row != existing {
		t.Errorf("expected row to be returned, got %+v", row)
	}
	if msg.Type != "reminder" {
		t.Errorf("expected message type 'reminder', got %s", msg.Type)
	}
}

func TestHandleProjectNotification_AllNotificationsSent(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
	s3Bucket := "test-bucket"
	existing := &models.ProjectNotification{HasSentWelcome: true, HasSentReminder: true}

	_, _, err := HandleProjectNotification(s3Bucket, project, existing)
	if err == nil || !strings.Contains(err.Error(), "AllNotificationsAlreadySent") {
		t.Errorf("expected 'AllNotificationsAlreadySent' error, got %v", err)
	}
}

func TestHandleProjectNotification_EmptyBucket(t *testing.T) {
	project := models.Project{Name: "TestProject", Date: "2025-11-11"}

	_, _, err := HandleProjectNotification("", project, nil)
	if err == nil || !strings.Contains(err.Error(), "S3 bucket name is required") {
		t.Errorf("expected 'S3 bucket name is required' error, got %v", err)
	}
}
