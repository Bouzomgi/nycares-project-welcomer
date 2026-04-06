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

type UpcomingProjectsService interface {
	GetUpcomingProjects(ctx context.Context, userSFID string) ([]domain.Project, error)
	SetCookies(cookies []*http.Cookie) error
}

func (s *HttpService) GetUpcomingProjects(ctx context.Context, userSFID string) ([]domain.Project, error) {
	if userSFID == "" {
		return nil, fmt.Errorf("userSFID is required")
	}

	req, err := s.buildUpcomingProjectsRequest(userSFID)
	if err != nil {
		return nil, fmt.Errorf("failed to build upcoming projects request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("upcoming projects request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return nil, fmt.Errorf("upcoming projects request failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var upcomingResp []dto.UpcomingResponse
	if err := json.Unmarshal(body, &upcomingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal upcoming projects response: %w", err)
	}

	if len(upcomingResp) == 0 {
		return nil, fmt.Errorf("upcoming projects response was empty")
	}

	projects, err := upcomingResp[0].ToDomainProjects()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *HttpService) buildUpcomingProjectsRequest(userSFID string) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/%s/1", endpoints.JoinPaths(s.baseUrl, endpoints.GetUpcomingProjectsPath), userSFID)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
