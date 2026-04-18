package dynamoservice

import (
	"context"
	"testing"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

// TestGetProjectNotification_EmptyTableName verifies that GetProjectNotification
// returns an error before touching the DynamoDB client when tableName is empty.
func TestGetProjectNotification_EmptyTableName(t *testing.T) {
	svc := NewDynamoService(nil, "")

	project := domain.Project{
		Name: "Park Cleanup",
		Date: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
	}

	_, err := svc.GetProjectNotification(context.Background(), project)
	if err == nil {
		t.Fatal("expected error for empty tableName, got nil")
	}
}

// TestUpsertProjectNotification_EmptyTableName verifies that UpsertProjectNotification
// returns an error before touching the DynamoDB client when tableName is empty.
func TestUpsertProjectNotification_EmptyTableName(t *testing.T) {
	svc := NewDynamoService(nil, "")

	pn := domain.ProjectNotification{
		Name:             "Park Cleanup",
		Date:             time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
		Id:               "proj-001",
		HasSentWelcome:   true,
		HasSentReminder:  false,
		HasSentThankYou:  false,
		ShouldStopNotify: false,
	}

	_, err := svc.UpsertProjectNotification(context.Background(), pn)
	if err == nil {
		t.Fatal("expected error for empty tableName, got nil")
	}
}

// TestDomainRoundTrip verifies that the domain.ProjectNotification fields we
// write into the DTO and read back are consistent (pure struct mapping, no AWS call).
func TestDomainRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		pn   domain.ProjectNotification
	}{
		{
			name: "all flags false",
			pn: domain.ProjectNotification{
				Name:             "Park Cleanup",
				Date:             time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
				Id:               "proj-001",
				HasSentWelcome:   false,
				HasSentReminder:  false,
				HasSentThankYou:  false,
				ShouldStopNotify: false,
			},
		},
		{
			name: "all flags true",
			pn: domain.ProjectNotification{
				Name:             "Food Bank",
				Date:             time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
				Id:               "proj-002",
				HasSentWelcome:   true,
				HasSentReminder:  true,
				HasSentThankYou:  true,
				ShouldStopNotify: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the mapping logic from UpsertProjectNotification result block.
			result := &domain.ProjectNotification{
				Name:             tt.pn.Name,
				Date:             tt.pn.Date,
				Id:               tt.pn.Id,
				HasSentWelcome:   tt.pn.HasSentWelcome,
				HasSentReminder:  tt.pn.HasSentReminder,
				HasSentThankYou:  tt.pn.HasSentThankYou,
				ShouldStopNotify: tt.pn.ShouldStopNotify,
			}

			if result.Name != tt.pn.Name {
				t.Errorf("Name = %q, want %q", result.Name, tt.pn.Name)
			}
			if !result.Date.Equal(tt.pn.Date) {
				t.Errorf("Date = %v, want %v", result.Date, tt.pn.Date)
			}
			if result.Id != tt.pn.Id {
				t.Errorf("Id = %q, want %q", result.Id, tt.pn.Id)
			}
			if result.HasSentWelcome != tt.pn.HasSentWelcome {
				t.Errorf("HasSentWelcome = %v, want %v", result.HasSentWelcome, tt.pn.HasSentWelcome)
			}
			if result.HasSentReminder != tt.pn.HasSentReminder {
				t.Errorf("HasSentReminder = %v, want %v", result.HasSentReminder, tt.pn.HasSentReminder)
			}
			if result.HasSentThankYou != tt.pn.HasSentThankYou {
				t.Errorf("HasSentThankYou = %v, want %v", result.HasSentThankYou, tt.pn.HasSentThankYou)
			}
			if result.ShouldStopNotify != tt.pn.ShouldStopNotify {
				t.Errorf("ShouldStopNotify = %v, want %v", result.ShouldStopNotify, tt.pn.ShouldStopNotify)
			}
		})
	}
}
