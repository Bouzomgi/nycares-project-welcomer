package httpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/httpclient"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

func (s *HttpService) GetSchedule(ctx context.Context, internalID string) ([]models.DetailedProject, error) {
	if internalID == "" {
		return nil, fmt.Errorf("internalID is required")
	}

	req, err := s.buildScheduleRequest(internalID)
	if err != nil {
		return nil, fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.client.SendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("schedule request failed: %w", err)
	}

	if err := httpclient.CheckResponse(resp); err != nil {
		if err != nil {
			return nil, fmt.Errorf("schedule request failed: %w", err)
		}
	}

	body, err := s.client.ReadBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var respArray []models.ScheduleResponse
	if err := json.Unmarshal(body, &respArray); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	if len(respArray) == 0 {
		return nil, fmt.Errorf("empty schedule response")
	}

	schedule := respArray[0]
	return flattenScheduleList(schedule), nil
}

func (s *HttpService) buildScheduleRequest(internalId string) (*http.Request, error) {
	getScheduleBaseUrl := s.baseURL + endpoints.GetSchedulePath
	urlStr := fmt.Sprintf("%s/%s", getScheduleBaseUrl, internalId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	return req, nil
}

// flattenScheduleList converts a ScheduleResponse into a slice of DetailedProjects
func flattenScheduleList(resp models.ScheduleResponse) []models.DetailedProject {
	var result []models.DetailedProject
	for _, project := range resp.Data.ScheduleList {
		result = append(result, project)
	}
	return result
}

// // GetProjectByNameAndDate finds a specific project by name and date
// func (s *HttpService) GetProjectByNameAndDate(projects []models.DetailedProject, name string, date time.Time) (*models.DetailedProject, error) {
// 	for _, project := range projects {
// 		if project.WebTitleFF == name && project.StartDate.Equal(date) {
// 			return &project, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("project not found: %s on %s", name, date.Format("2006-01-02"))
// }
