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
	GetTodayProjects(ctx context.Context, userSFID string) ([]domain.Project, error)
	SetCookies(cookies []*http.Cookie) error
}

func (s *HttpService) GetUpcomingProjects(ctx context.Context, userSFID string) ([]domain.Project, error) {
	return s.fetchProjects(ctx, userSFID, endpoints.GetUpcomingProjectsPath, "upcoming")
}

func (s *HttpService) GetTodayProjects(ctx context.Context, userSFID string) ([]domain.Project, error) {
	return s.fetchProjects(ctx, userSFID, endpoints.GetTodayProjectsPath, "today")
}

func (s *HttpService) fetchProjects(ctx context.Context, userSFID, path, label string) ([]domain.Project, error) {
	if userSFID == "" {
		return nil, fmt.Errorf("userSFID is required")
	}

	req, err := s.buildProjectsRequest(userSFID, path)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s projects request: %w", label, err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s projects request failed: %w", label, err)
	}

	if err := CheckResponse(resp); err != nil {
		return nil, fmt.Errorf("%s projects request failed: %w", label, err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s projects response body: %w", label, err)
	}

	var upcomingResp []dto.UpcomingResponse
	if err := json.Unmarshal(body, &upcomingResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s projects response: %w", label, err)
	}

	if len(upcomingResp) == 0 {
		return []domain.Project{}, nil
	}

	projects, err := upcomingResp[0].ToDomainProjects()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *HttpService) buildProjectsRequest(userSFID, path string) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s/%s/1", endpoints.JoinPaths(s.baseUrl, path), userSFID)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
