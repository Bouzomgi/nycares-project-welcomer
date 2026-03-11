package sesservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESService struct {
	client    *sesv2.Client
	sender    string
	recipient string
}

func NewSESService(client *sesv2.Client, sender, recipient string) *SESService {
	return &SESService{client: client, sender: sender, recipient: recipient}
}

func (s *SESService) SendHTMLEmail(ctx context.Context, subject, htmlBody, plainText string) error {
	_, err := s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.sender),
		Destination: &types.Destination{
			ToAddresses: []string{s.recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Html: &types.Content{Data: aws.String(htmlBody)},
					Text: &types.Content{Data: aws.String(plainText)},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send HTML email: %w", err)
	}
	return nil
}

func (s *SESService) SendPlainEmail(ctx context.Context, subject, body string) error {
	_, err := s.client.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(s.sender),
		Destination: &types.Destination{
			ToAddresses: []string{s.recipient},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Text: &types.Content{Data: aws.String(body)},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to send plain email: %w", err)
	}
	return nil
}
