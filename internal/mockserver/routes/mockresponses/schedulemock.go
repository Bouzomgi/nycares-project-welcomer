package mockresponses

import (
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

func MockScheduleResponse() dto.ScheduleResponse {

	return dto.ScheduleResponse{
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
			ScheduleList: map[string]dto.Project{
				"1": {
					Role:               "Volunteer",
					FamilyFriendlyRole: nil,
					Id:                 "a1Bxx0000001XYZ",
					Status:             "Scheduled",
					WebTitleFF:         "Test Project",
					StartDate:          time.Now().Format("2006-01-02"),
					ActivityStartTime:  "09:00",
					EndDate:            time.Now().Add(2 * time.Hour).Format("2006-01-02"),
					ActivityEndTime:    "12:00",
					CampaignId:         "11111111-1111-1111-1111-111111111111",
				},
			},
			UpcomingCount:        1,
			PlusCount:            0,
			ShowNewFunctionality: true,
		},
	}
}
