package models

import (
	"fmt"
	"time"
)

type DetailedProject struct {
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
	ScheduleList         map[string]DetailedProject `json:"ScheduleList__tl"`
	UpcomingCount        int                        `json:"UpcomingCount__tl"`
	PlusCount            int                        `json:"PlusCount__tl"`
	ShowNewFunctionality bool                       `json:"ShowNewFunctionality__tl"`
}

type ScheduleResponse struct {
	Success          bool         `json:"success"`
	Data             ScheduleData `json:"data"`
	Message          string       `json:"message"`
	Command          string       `json:"command"`
	IsUserTeamLeader bool         `json:"is_user_team_leader"`
	UserSFID         string       `json:"user_sf_id"`
}

type Project struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func reduceToProject(DetailedProject DetailedProject) Project {
	return Project{
		Name: DetailedProject.WebTitleFF,
		Date: DetailedProject.StartDate,
	}
}

func ReduceToProjects(detailed []DetailedProject) []Project {
	projects := make([]Project, 0, len(detailed))
	for _, dp := range detailed {
		projects = append(projects, reduceToProject(dp))
	}
	return projects
}

func FilterUpcomingProjects(projects []DetailedProject, within time.Duration) ([]DetailedProject, error) {
	now := time.Now()
	var upcoming []DetailedProject

	for _, project := range projects {
		projectStart, err := time.Parse("2006-01-02", project.StartDateTime)
		if err != nil {
			return nil, fmt.Errorf("could not parse project start %s: %w", project.StartDateTime, err)
		}

		if projectStart.Sub(now) <= within {
			upcoming = append(upcoming, project)
		}
	}

	return upcoming, nil
}
