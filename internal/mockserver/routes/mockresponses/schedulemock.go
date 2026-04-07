package mockresponses

import (
	"os"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

type ProjectConfig struct {
	Name string
	Date time.Time
	Id   string
}

func currentDate() time.Time {
	if dateStr := os.Getenv("NYCARES_CURRENT_DATE"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			return t
		}
	}
	return time.Now()
}

func MockUpcomingResponse(projects []ProjectConfig) []dto.UpcomingResponse {
	if projects == nil {
		now := currentDate()
		projects = []ProjectConfig{
			{Name: "Test Project", Date: now.AddDate(0, 0, 6), Id: "MOCKSESSION0000001"},
		}
	}

	sessions := make([]dto.UpcomingSession, 0, len(projects))
	for _, p := range projects {
		sessions = append(sessions, dto.UpcomingSession{
			Name:               p.Name,
			PublicSessionName:  p.Name,
			FamilyFriendlyRole: nil,
			SessionID:          p.Id,
			Status:             "Published",
			SessionStartDate:   p.Date.Format("2006-01-02"),
			SessionStartTime:   "10:00:00.000Z",
			SessionEndDate:     p.Date.Format("2006-01-02"),
			SessionEndTime:     "12:00:00.000Z",
			DatetimeState:      "upcoming",
			AWSChimeChannelID:  utils.NewUUID(),
			RegistrationStatus: "signed up",
			IsTeamLeader:       true,
		})
	}

	return []dto.UpcomingResponse{
		{
			Success:          true,
			Data:             sessions,
			Message:          "Retrieved dashboard upcoming campaign(s).",
			Page:             "1",
			Command:          "SessionActiveUpcoming",
			IsUserTeamLeader: true,
			UserSFID:         "MOCKSFID0000000001A",
			IsUserFlagged: dto.UserFlagged{
				Deactivated: false,
			},
			IsVolunteer:            false,
			UserFamilyFriendlyRole: nil,
			UserAWSID:              "MOCKAWSID0000000001",
			OrientationURL:         "/sites/default/files/trainings/volunteer-orientation/index.html",
			VIFURL:                 "/volunteer-information-form",
		},
	}
}
