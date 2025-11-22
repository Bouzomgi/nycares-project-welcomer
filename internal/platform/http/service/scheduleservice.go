package httpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

type ScheduleService interface {
	GetSchedule(ctx context.Context, internalID string) ([]domain.Project, error)
	SetCookies(cookies []*http.Cookie) error
}

func (s *HttpService) GetSchedule(ctx context.Context, internalID string) ([]domain.Project, error) {
	if internalID == "" {
		return nil, fmt.Errorf("internalID is required")
	}

	req, err := s.buildScheduleRequest(internalID)
	if err != nil {
		return nil, fmt.Errorf("failed to build schedule request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("schedule request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return nil, fmt.Errorf("schedule request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var schedResp []dto.ScheduleResponse
	if err := json.Unmarshal(body, &schedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	projects, err := schedResp[0].ToDomainProjects()
	if err != nil {
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("empty schedule response")
	}

	return projects, nil
}

func (s *HttpService) buildScheduleRequest(internalId string) (*http.Request, error) {
	getScheduleBaseUrl := endpoints.JoinPaths(endpoints.BaseUrl, endpoints.GetSchedulePath)
	urlStr := fmt.Sprintf("%s/%s", getScheduleBaseUrl, internalId)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
