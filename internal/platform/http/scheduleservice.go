package httpservice

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type project struct {
	Role               string  `json:"Role__c"`
	FamilyFriendlyRole *string `json:"Family_Friendly_Role__c"`
	Id                 string  `json:"Id"`
	Status             string  `json:"Status"`
	WebTitleFF         string  `json:"Web_Title_FF__c"`
	StartDate          string  `json:"StartDate"`
	ActivityStartTime  string  `json:"Activity_Start_Time__c"`
	EndDate            string  `json:"EndDate"`
	ActivityEndTime    string  `json:"Activity_End_Time__c"`
	ContactArray       []struct {
		Id                 string `json:"Id"`
		Name               string `json:"Name"`
		Role               string `json:"Role__c"`
		IsTeamLeader       bool   `json:"IsTeamLeader__tl"`
		DisplayNameAndRole string `json:"DisplayNameAndRole__tl"`
	} `json:"ContactArray__tl"`
	CampaignId         string `json:"CampaignId__tl"`
	StartDateTime      string `json:"StartDateTime__tl"`
	EndDateTime        string `json:"EndDateTime__tl"`
	ContactDisplayList string `json:"ContactDisplayList__tl"`
	DayOfWeek          string `json:"DayOfWeek__tl"`
}

type scheduleData struct {
	ScheduleList         map[string]project `json:"ScheduleList__tl"`
	UpcomingCount        int                `json:"UpcomingCount__tl"`
	PlusCount            int                `json:"PlusCount__tl"`
	ShowNewFunctionality bool               `json:"ShowNewFunctionality__tl"`
}

type userFlagged struct {
	Deactivated       bool `json:"Deactivated"`
	DoubleDeactivated bool `json:"DoubleDeactivated"`
	NeedsOrientation  bool `json:"NeedsOrientation"`
	NeedsVIF          bool `json:"NeedsVIF"`
	YesToConviction   bool `json:"YesToConviction"`
	Ineligible        bool `json:"Ineligible"`
	FamilyDeactivated bool `json:"FamilyDeactivated"`
}

type scheduleResponse struct {
	Success                bool         `json:"success"`
	Data                   scheduleData `json:"data"`
	Message                string       `json:"message"`
	Command                string       `json:"command"`
	IsUserTeamLeader       bool         `json:"is_user_team_leader"`
	UserSFID               string       `json:"user_sf_id"`
	IsUserFlagged          userFlagged  `json:"is_user_flagged"`
	UserFamilyFriendlyRole *string      `json:"user_family_friendly_role"`
	OrientationURL         string       `json:"orientation_url"`
	VIFURL                 string       `json:"vif_url"`
}

/////////

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

	var schedResp []scheduleResponse
	if err := json.Unmarshal(body, &schedResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schedule response: %w", err)
	}

	projects, err := schedResp[0].toDomainProjects()
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

	// TODO: can I remove this? And all others?
	req.Header.Set("User-Agent", "Mozilla/5.0")

	return req, nil
}

func (sr scheduleResponse) toDomainProjects() ([]domain.Project, error) {
	var projects []domain.Project
	for _, p := range sr.Data.ScheduleList {
		projectDate, err := utils.StringToDate(p.StartDate)

		if err != nil {
			return nil, fmt.Errorf("could not parse validate date from site's schedule response")
		}

		projects = append(projects, domain.Project{
			Name: p.WebTitleFF,
			Date: projectDate,
			Id:   p.Id,
		})
	}
	return projects, nil
}
