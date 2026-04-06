package mockresponses

import (
	"github.com/Bouzomgi/nycares-project-welcomer/internal/mockserver/utils"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/dto"
)

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
						{200, 0.123},
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
