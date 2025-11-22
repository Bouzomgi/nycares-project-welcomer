package dto

import (
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

/////// GetSchedule

type ScheduleResponse struct {
	Success                bool         `json:"success"`
	Data                   ScheduleData `json:"data"`
	Message                string       `json:"message"`
	Command                string       `json:"command"`
	IsUserTeamLeader       bool         `json:"is_user_team_leader"`
	UserSFID               string       `json:"user_sf_id"`
	IsUserFlagged          UserFlagged  `json:"is_user_flagged"`
	UserFamilyFriendlyRole *string      `json:"user_family_friendly_role"`
	OrientationURL         string       `json:"orientation_url"`
	VIFURL                 string       `json:"vif_url"`
}

type Project struct {
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

type ScheduleData struct {
	ScheduleList         map[string]Project `json:"ScheduleList__tl"`
	UpcomingCount        int                `json:"UpcomingCount__tl"`
	PlusCount            int                `json:"PlusCount__tl"`
	ShowNewFunctionality bool               `json:"ShowNewFunctionality__tl"`
}

type UserFlagged struct {
	Deactivated       bool `json:"Deactivated"`
	DoubleDeactivated bool `json:"DoubleDeactivated"`
	NeedsOrientation  bool `json:"NeedsOrientation"`
	NeedsVIF          bool `json:"NeedsVIF"`
	YesToConviction   bool `json:"YesToConviction"`
	Ineligible        bool `json:"Ineligible"`
	FamilyDeactivated bool `json:"FamilyDeactivated"`
}

func (sr ScheduleResponse) ToDomainProjects() ([]domain.Project, error) {
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
