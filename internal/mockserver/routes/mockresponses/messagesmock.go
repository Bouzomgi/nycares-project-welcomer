package mockresponses

import (
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

func MockCampaignResponse() dto.CampaignResponse {
	return dto.CampaignResponse{
		{
			Campaign: dto.Campaign{
				Name: "Mock Project A",
				RecordType: dto.RecordType{
					Attributes: dto.RecordTypeAttributes{
						Type: "RecordType",
						URL:  "/services/data/vXX.X/sobjects/RecordType/012ABC",
					},
					Name: "Volunteer Event",
				},
				ProgramType:             "Environmental",
				Status:                  "Active",
				WebTitleFF:              "Mock Family-Friendly Title",
				StartDate:               "2025-11-05",
				ParentId:                "001XYZ",
				ActivityStartTime:       "09:00",
				EndDate:                 "2025-11-05",
				ActivityEndTime:         "12:00",
				WebPublicationStartDate: "2025-10-01",
				WebPublicationEndDate:   "2025-12-31",
				SpecialProject:          "None",
				CommittedProjectDateRange: func() *string {
					s := "Nov 5 - Nov 7"
					return &s
				}(),
				Borough:            "Brooklyn",
				WebsiteAddress:     "https://example.org/project",
				TeamLeaderToolsId:  nil,
				DrupalID:           nil,
				AttendanceTaken:    false,
				ProjectDescription: "A mock project for testing.",
				TeamLeaderNotes:    nil,
				ProjectLogistics:   nil,
				RepOnSite: dto.Contact{
					Attributes: dto.RecordTypeAttributes{
						Type: "Contact",
						URL:  "/services/data/vXX.X/sobjects/Contact/003AAA",
					},
					Name: "John Rep",
				},
				CommunityPartnerName: "Mock Community Partner",
				Directions:           "Mock directions",
				FullCapacity:         "30",
				AWSChimeChannelArn:   "arn:aws:chime:us-east-1:123456789012:channel/mock",
				TeamLeaderContact:    "jane.doe@example.com",
				TeamLeadersList:      nil,
				PinnedChatMessage:    nil,
				NumOfRegistration:    10,
				CapacityRemaining:    20,
				Agency: dto.Agency{
					Attributes: dto.RecordTypeAttributes{
						Type: "Agency",
						URL:  "/services/data/vXX.X/sobjects/Agency/00A111",
					},
					Name:        "Mock Agency",
					Description: "A mock agency for unit tests.",
				},
				AgencyDescription:             nil,
				GeneralInterestCampaign:       false,
				AttendanceSignedUpCount:       5,
				DatetimeState:                 "Upcoming",
				OrientationVIFNotRequired:     true,
				HumanReadableDate:             "November 5, 2025",
				Id:                            "701ABCDEF",
				RegistrationId:                "reg-123",
				UserStatus:                    "Registered",
				UserRole:                      "Volunteer",
				RecordTypeTL:                  "Event",
				IsTeamLeader:                  false,
				SpecialProjectTL:              []string{},
				IsMultiSession:                false,
				IsFirstSession:                true,
				IsCommittedProject:            false,
				CommittedProjectDateRangeTL:   []string{},
				IsTeenFriendly:                true,
				IsFamilyFriendly:              true,
				StartDateTimeTL:               "2025-11-05T09:00:00Z",
				EndDateTimeTL:                 "2025-11-05T12:00:00Z",
				ProjectOccurrenceState:        "Future",
				ProjectIsUpcoming:             true,
				CommunityPartnerTL:            "Mock CP",
				SiteLocation:                  "123 Test St",
				RepOnSiteName:                 "John Rep",
				SiteAddressTL:                 "123 Test St, Brooklyn, NY",
				AttendanceTakenTL:             false,
				SiteDescription:               "Mock site description",
				CommunityPartnerNameTL:        "Mock CP",
				CommunityPartnerDescriptionTL: "Mock CP Description",
				AWSChimeChannelID:             "mock-channel-id",
				CurrentUserID:                 "005USERXYZ",
				ProjectManagerID:              "005PM12345",
				ProjectManagerName:            "PM Name",
				ProjectManagerChannelLink:     "https://mock-channel-link",
				TeamLeaderID:                  "005TL12345",
				TeamLeaderName:                "Team Leader",
				TeamLeaderFirstName:           "Leader",
				TeamLeaderChannelLink:         "https://mock-tl-link",
				Bookmarked:                    false,
				IsRecent:                      false,
				CampaignSiblings:              []dto.Campaign{},
			},
		},
	}
}

///////

func MockSendMessageResponse() dto.SendMessageResponse {
	return dto.SendMessageResponse{
		Success: true,
		Message: "Message sent successfully",
		Data: dto.Data{
			ChannelArn: "arn:aws:chime:us-east-1:123456789012:app-instance/abcdefg/channel/12345678",
			MessageId:  utils.GenerateRandomMessageId(),
			Metadata: dto.Metadata{
				StatusCode:   200,
				EffectiveURI: "https://service.aws.com/sendMessage",
				Headers: dto.Headers{
					Date:           "Mon, 01 Jan 2025 00:00:00 GMT",
					ContentType:    "application/json",
					ContentLength:  "123",
					Connection:     "keep-alive",
					XAmznRequestID: "req-abcdef123456",
				},
				TransferStats: dto.TransferStats{
					HTTP: [][]interface{}{
						{200, 0.123}, // status code, time taken in seconds
					},
				},
			},
		},
	}
}

///////

func MockChannelMessagesResponse() dto.ChannelMessagesResponse {

	return dto.ChannelMessagesResponse{
		{
			ChannelMessages: []dto.ChannelMessage{
				{
					MessageId: utils.GenerateRandomMessageId(),
					Content:   "This is a mock message",
					Metadata: dto.MessageMetadata{
						SenderSFIDTL: "123456",
						Attachments:  []interface{}{},
					},
					Type:                 "STANDARD",
					CreatedTimestamp:     "2025-11-05T09:00:00Z",
					LastUpdatedTimestamp: "2025-11-05T09:00:00Z",
					Sender: dto.MessageSender{
						Arn:                  "arn:aws:chime:us-east-1:123456789012:app-instance/12345678",
						Name:                 "John Doe",
						RoleC:                "Participant",
						FirstNameLastInitial: "J.D.",
						PhotoURL:             "https://example.com/photo.jpg",
					},
					Redacted:             false,
					ChannelType:          "CHAT",
					ProjectIdTL:          "123456",
					JSONRepresentationTL: "{}",
				},
			},
		},
	}
}

////////

func MockPostPinResponse() dto.PostPinResponse {
	return dto.PostPinResponse{
		Scalar: true,
	}
}
