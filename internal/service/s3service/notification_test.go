package s3service

// import (
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
// )

// func TestComputeProjectMessage_NoExistingNotification(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
// 	s3Bucket := "test-bucket"

// 	service := NewS3Service(s3Bucket)
// 	row, msg, err := service.ComputeProjectMessage(project, nil)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if row != nil {
// 		t.Fatalf("expected row to be nil, got %+v", err)
// 	}
// 	if msg.Type != "welcome" {
// 		t.Fatalf("expected message type 'welcome', got %s", msg.Type)
// 	}
// }

// func TestComputeProjectMessage_ShouldStopNotify(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
// 	s3Bucket := "test-bucket"
// 	existing := &models.ProjectNotification{ShouldStopNotify: true}

// 	service := NewS3Service(s3Bucket)
// 	_, _, err := service.ComputeProjectMessage(project, existing)

// 	if err == nil || !strings.Contains(err.Error(), "notifications disabled") {
// 		t.Errorf("expected 'notifactions disabled' error, got %v", err)
// 	}
// }

// func TestComputeProjectMessage_WelcomeNotSent(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
// 	s3Bucket := "test-bucket"
// 	existing := &models.ProjectNotification{HasSentWelcome: true}

// 	service := NewS3Service(s3Bucket)
// 	row, msg, err := service.ComputeProjectMessage(project, existing)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if row != existing {
// 		t.Errorf("expected row to be returned, got %+v", row)
// 	}
// 	if msg.Type != "welcome" {
// 		t.Errorf("expected message type 'welcome', got %s", msg.Type)
// 	}
// }

// func TestComputeProjectMessage_ReminderNotSent(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
// 	s3Bucket := "test-bucket"
// 	existing := &models.ProjectNotification{HasSentWelcome: true, HasSentReminder: false}

// 	service := NewS3Service(s3Bucket)
// 	row, msg, err := service.ComputeProjectMessage(project, existing)

// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if row != existing {
// 		t.Errorf("expected row to be returned, got %+v", row)
// 	}
// 	if msg.Type != "welcome" {
// 		t.Errorf("expected message type 'reminder', got %s", msg.Type)
// 	}
// }

// func TestComputeProjectMessage_AllNotificationSent(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}
// 	s3Bucket := "test-bucket"
// 	existing := &models.ProjectNotification{HasSentWelcome: true, HasSentReminder: true}

// 	service := NewS3Service(s3Bucket)
// 	_, _, err := service.ComputeProjectMessage(project, existing)

// 	if err != nil || !strings.Contains(err.Error(), "all notifications already sent") {
// 		t.Fatalf("expected 'all notifications already sent' error, got %v", err)
// 	}
// }

// func TestComputeProjectMessage_EmptyBucket(t *testing.T) {
// 	project := models.Project{Name: "TestProject", Date: "2025-11-11"}

// 	service := NewS3Service("")
// 	_, _, err := service.ComputeProjectMessage(project, nil)

// 	if err != nil || !strings.Contains(err.Error(), "S3 bucket name is required") {
// 		t.Errorf("expected 'S3 bucket name is required' error, got %v", err)
// 	}
// }

// func TestShouldSendWelcome(t *testing.T) {
// 	// Test cases
// 	tests := []struct {
// 		name        string
// 		projectDate time.Time
// 		expected    bool
// 	}{
// 		{
// 			name:        "Should send - 3 days before",
// 			projectDate: time.Now().AddDate(0, 0, 3),
// 			expected:    true,
// 		},
// 		{
// 			name:        "Should send - 10 days before",
// 			projectDate: time.Now().AddDate(0, 0, 10),
// 			expected:    false,
// 		},
// 		{
// 			name:        "Should send - 1 day after",
// 			projectDate: time.Now().AddDate(0, 0, -1),
// 			expected:    false,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := ShouldSendWelcome(tc.projectDate)
// 			if result != tc.expected {
// 				t.Errorf("ShouldSendWecome(%v) = %v; want %v", tc.projectDate, result, tc.expected)
// 			}
// 		})
// 	}
// }

// func TestShouldSendReminder(t *testing.T) {
// 	// Test cases
// 	tests := []struct {
// 		name        string
// 		projectDate time.Time
// 		expected    bool
// 	}{
// 		{
// 			name:        "Should send - 1 day before",
// 			projectDate: time.Now().AddDate(0, 0, 1),
// 			expected:    true,
// 		},
// 		{
// 			name:        "Should send - 5 days before",
// 			projectDate: time.Now().AddDate(0, 0, 5),
// 			expected:    false,
// 		},
// 		{
// 			name:        "Should send - 1 day after",
// 			projectDate: time.Now().AddDate(0, 0, -1),
// 			expected:    false,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			result := ShouldSendReminder(tc.projectDate)
// 			if result != tc.expected {
// 				t.Errorf("ShouldSendReminder(%v) = %v; want %v", tc.projectDate, result, tc.expected)
// 			}
// 		})
// 	}
// }
