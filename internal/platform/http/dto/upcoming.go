package dto

import (
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type UserFlagged struct {
	Deactivated       bool `json:"Deactivated"`
	DoubleDeactivated bool `json:"DoubleDeactivated"`
	NeedsOrientation  bool `json:"NeedsOrientation"`
	NeedsVIF          bool `json:"NeedsVIF"`
	YesToConviction   bool `json:"YesToConviction"`
	Ineligible        bool `json:"Ineligible"`
	FamilyDeactivated bool `json:"FamilyDeactivated"`
}

type UpcomingSession struct {
	Name               string  `json:"Name"`
	PublicSessionName  string  `json:"Public_Session_Name__c"`
	FamilyFriendlyRole *string `json:"Family_Friendly_Role__c"`
	SessionID          string  `json:"Session__c"`
	Status             string  `json:"Status__c"`
	SessionStartDate   string  `json:"Session_Start_Date__c"`
	SessionStartTime   string  `json:"Session_Start_Time__c"`
	SessionEndDate     string  `json:"Session_End_Date__c"`
	SessionEndTime     string  `json:"Session_End_Time__c"`
	DatetimeState      string  `json:"Datetime_State__c"`
	AWSChimeChannelID  string  `json:"AWS_Chime_Channel_Arn_Channel_Id__c"`
	RegistrationStatus string  `json:"Registration_Status__tl"`
	IsTeamLeader       bool    `json:"IsTeamLeader__tl"`
}

type UpcomingResponse struct {
	Success                bool              `json:"success"`
	Data                   []UpcomingSession `json:"data"`
	Message                string            `json:"message"`
	Page                   string            `json:"page"`
	Command                string            `json:"command"`
	IsUserTeamLeader       bool              `json:"is_user_team_leader"`
	UserSFID               string            `json:"user_sf_id"`
	IsUserFlagged          UserFlagged       `json:"is_user_flagged"`
	IsVolunteer            bool              `json:"is_volunteer"`
	UserFamilyFriendlyRole *string           `json:"user_family_friendly_role"`
	UserAWSID              string            `json:"user_aws_id"`
	OrientationURL         string            `json:"orientation_url"`
	VIFURL                 string            `json:"vif_url"`
}

func (ur UpcomingResponse) ToDomainProjects() ([]domain.Project, error) {
	var projects []domain.Project
	for _, s := range ur.Data {
		projectDate, err := utils.StringToDate(s.SessionStartDate)
		if err != nil {
			return nil, fmt.Errorf("could not parse date from upcoming projects response")
		}
		projects = append(projects, domain.Project{
			Name:      s.PublicSessionName,
			Date:      projectDate,
			Id:        s.SessionID,
			ChannelId: s.AWSChimeChannelID,
		})
	}
	return projects, nil
}
