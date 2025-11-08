package main

// MockScheduleResponse returns a mock ScheduleResponse for testing
func MockScheduleResponse() ScheduleResponse {
	familyFriendly := "Yes"

	return ScheduleResponse{
		Success:          true,
		Message:          "Fetched schedule successfully",
		Command:          "GetSchedule",
		IsUserTeamLeader: false,
		UserSFID:         "003ABC123456XYZ",
		Data: ScheduleData{
			UpcomingCount:        2,
			PlusCount:            1,
			ShowNewFunctionality: true,
			ScheduleList: map[string]CompleteProject{
				"a0B1X00000Example1": {
					Role:               "Volunteer",
					FamilyFriendlyRole: &familyFriendly,
					Id:                 "a0B1X00000Example1",
					Status:             "Confirmed",
					WebTitleFF:         "Community Garden Cleanup",
					StartDate:          "2025-11-06",
					ActivityStartTime:  "09:00",
					EndDate:            "2025-11-06",
					ActivityEndTime:    "12:00",
					CampaignId:         "7011X000000ABCD",
					StartDateTime:      "2025-11-06T09:00:00Z",
					EndDateTime:        "2025-11-06T12:00:00Z",
					ContactDisplayList: "John Doe (Team Leader), Jane Smith (Volunteer)",
					DayOfWeek:          "Thursday",
					ContactArray: []struct {
						Id                 string `json:"Id"`
						Name               string `json:"Name"`
						Role               string `json:"Role__c"`
						IsTeamLeader       bool   `json:"IsTeamLeader__tl"`
						DisplayNameAndRole string `json:"DisplayNameAndRole__tl"`
					}{
						{
							Id:                 "0031X00000Leader1",
							Name:               "John Doe",
							Role:               "Team Leader",
							IsTeamLeader:       true,
							DisplayNameAndRole: "John Doe (Team Leader)",
						},
						{
							Id:                 "0031X00000Volunteer1",
							Name:               "Jane Smith",
							Role:               "Volunteer",
							IsTeamLeader:       false,
							DisplayNameAndRole: "Jane Smith (Volunteer)",
						},
					},
				},
				"a0B1X00000Example2": {
					Role:               "Volunteer",
					FamilyFriendlyRole: nil,
					Id:                 "a0B1X00000Example2",
					Status:             "Pending",
					WebTitleFF:         "Food Pantry Assistance",
					StartDate:          "2025-11-10",
					ActivityStartTime:  "13:00",
					EndDate:            "2025-11-10",
					ActivityEndTime:    "16:00",
					CampaignId:         "7011X000000EFGH",
					StartDateTime:      "2025-11-10T13:00:00Z",
					EndDateTime:        "2025-11-10T16:00:00Z",
					ContactDisplayList: "Mary Johnson (Coordinator)",
					DayOfWeek:          "Monday",
					ContactArray: []struct {
						Id                 string `json:"Id"`
						Name               string `json:"Name"`
						Role               string `json:"Role__c"`
						IsTeamLeader       bool   `json:"IsTeamLeader__tl"`
						DisplayNameAndRole string `json:"DisplayNameAndRole__tl"`
					}{
						{
							Id:                 "0031X00000Coordinator1",
							Name:               "Mary Johnson",
							Role:               "Coordinator",
							IsTeamLeader:       false,
							DisplayNameAndRole: "Mary Johnson (Coordinator)",
						},
					},
				},
			},
		},
	}
}
