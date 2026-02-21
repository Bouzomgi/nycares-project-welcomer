package mockresponses

import (
	"fmt"
	"os"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

type ProjectConfig struct {
	Name       string
	Date       time.Time
	Id         string
	CampaignId string
}

func currentDate() time.Time {
	if dateStr := os.Getenv("NYCARES_CURRENT_DATE"); dateStr != "" {
		if t, err := time.Parse("2006-01-02", dateStr); err == nil {
			return t
		}
	}
	return time.Now()
}

func MockScheduleResponse(projects []ProjectConfig) []dto.ScheduleResponse {
	scheduleList := make(map[string]dto.Project)

	if projects == nil {
		now := currentDate()
		projects = []ProjectConfig{
			{Name: "Test Project", Date: now, Id: "a1Bxx0000001XYZ", CampaignId: "11111111-1111-1111-1111-111111111111"},
		}
	}

	for i, p := range projects {
		key := fmt.Sprintf("%d", i+1)
		scheduleList[key] = dto.Project{
			Role:               "Volunteer",
			FamilyFriendlyRole: nil,
			Id:                 p.Id,
			Status:             "Scheduled",
			WebTitleFF:         p.Name,
			StartDate:          p.Date.Format("2006-01-02"),
			ActivityStartTime:  "09:00",
			EndDate:            p.Date.Format("2006-01-02"),
			ActivityEndTime:    "12:00",
			CampaignId:         p.CampaignId,
		}
	}

	return []dto.ScheduleResponse{
		{
			Success:          true,
			Message:          "Mock schedule",
			Command:          "GetSchedule",
			IsUserTeamLeader: false,
			UserSFID:         "005xx000001Sv6dAAC",
			IsUserFlagged: dto.UserFlagged{
				Deactivated: false,
			},
			UserFamilyFriendlyRole: nil,
			OrientationURL:         "https://example.com/orientation",
			VIFURL:                 "https://example.com/vif",
			Data: dto.ScheduleData{
				ScheduleList:         scheduleList,
				UpcomingCount:        len(scheduleList),
				PlusCount:            0,
				ShowNewFunctionality: true,
			},
		},
	}
}
