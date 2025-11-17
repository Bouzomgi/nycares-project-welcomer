package httpservice

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

func setupMockScheduleServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Check if the path contains the expected endpoint
		if r.URL.Path != "/api/schedule/retrieve/test-id" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Create a mock schedule response
		now := time.Now()
		tomorrow := now.AddDate(0, 0, 1)

		mockProject := models.DetailedProject{
			Id:         "123",
			WebTitleFF: "Test Project",
			Status:     "Active",
			StartDate:  tomorrow.String(),
			EndDate:    tomorrow.Add(2 * time.Hour).String(),
		}

		mockScheduleData := models.ScheduleData{
			ScheduleList: map[string]models.DetailedProject{
				"123": mockProject,
			},
			UpcomingCount: 1,
		}

		mockResponse := []models.ScheduleResponse{
			{
				Success: true,
				Data:    mockScheduleData,
				Message: "Success",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
}

func TestGetSchedule(t *testing.T) {
	server := setupMockScheduleServer()
	defer server.Close()

	mockHttpService, err := NewHttpService(server.URL)
	if err != nil {
		t.Fatalf("Could not create mock server: %v", err)
	}

	// Test GetSchedule
	ctx := context.Background()
	projects, err := mockHttpService.GetSchedule(ctx, "test-id")

	// Check results
	if err != nil {
		t.Fatalf("Expect GetSchedule to succeed, but got error: %v", err)
	}

	if len(projects) != 1 {
		t.Fatalf("Expected 1 project, got %d", len(projects))
	}

	project := projects[0]
	if project.WebTitleFF != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got %s", project.WebTitleFF)
	}
}

// func TestGetProjectByNameAndDate(t *testing.T) {
// 	// Create test projects
// 	now := time.Now()
// 	tomorrow := now.AddDate(0, 0, 1)

// 	projects := []models.DetailedProject{
// 		{
// 			Id:         "123",
// 			WebTitleFF: "Project A",
// 			Status:     "Active",
// 			StartDate:  tomorrow.String(),
// 		},
// 		{
// 			Id:         "456",
// 			WebTitleFF: "Project B",
// 			Status:     "Active",
// 			StartDate:  now.AddDate(0, 0, 2).String(),
// 		},
// 	}

// 	server := setupMockScheduleServer()
// 	defer server.Close()

// 	mockHttpService, err := NewHttpService(server.URL)
// 	if err != nil {
// 		t.Fatalf("Could not create mock server: %v", err)
// 	}

// 	// Test finding existing project
// 	project, err := projectService.GetProjectByNameAndDate(projects, "Project A", tomorrow)
// 	if err != nil {
// 		t.Fatalf("Expected to find project, but got error %v", err)
// 	}

// 	if project.ID != "123" {
// 		t.Errorf("Expected project ID '123', got %s", project.ID)
// 	}

// 	// Test finding non-existent project
// 	_, err = projectService.GetProjectByNameAndDate(projects, "Project C", tomorrow)
// 	if err == nil {
// 		t.Error("Expected error for non-existent project, but got nil")
// 	}
// }

func TestFilterUpcomingProjects(t *testing.T) {
	// Create test projects
	now := time.Now()

	projects := []models.DetailedProject{
		{
			Id:            "123",
			WebTitleFF:    "Soon Project",
			StartDateTime: now.Add(24 * time.Hour).Format("2006-01-02"),
		},
		{
			Id:            "123",
			WebTitleFF:    "Far Project",
			StartDateTime: now.Add(10 * 24 * time.Hour).Format("2006-01-02"),
		},
	}

	// Test filter in projects within 3 days
	upcoming, error := models.FilterUpcomingProjects(projects, 3*24*time.Hour)
	if error != nil {
		t.Fatalf("Error filtering project list: %s", error)
	}

	if len(upcoming) != 1 {
		t.Fatalf("Expected 1 upcoming project, got %d", len(upcoming))
	}

	if upcoming[0].Id != "123" {
		t.Fatalf("Expected project ID '123', got %s", upcoming[0].Id)
	}

	// Test filter in projects within 2 weeks
	upcoming, error = models.FilterUpcomingProjects(projects, 14*24*time.Hour)
	if error != nil {
		t.Fatalf("Error filtering project list: %s", error)
	}
	if len(upcoming) != 2 {
		t.Fatalf("Expected 2 upcoming projects, got %d", len(upcoming))
	}
}
