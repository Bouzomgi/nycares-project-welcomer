package snsservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type NotificationService interface {
	PublishHTMLEmailNotification(ctx context.Context, plainText, htmlBody, subject string) (string, error)
}

// PublishHTMLEmailNotification publishes an HTML email notification to the SNS topic.
// The message body is a JSON object with "htmlBody" and "plainText" fields.
// A "format: html" message attribute signals the SES forwarder Lambda to send HTML email.
func (s *SNSService) PublishHTMLEmailNotification(ctx context.Context, plainText, htmlBody, subject string) (string, error) {
	if s.topicARN == "" {
		return "", fmt.Errorf("topicARN cannot be empty")
	}
	if subject == "" {
		return "", fmt.Errorf("subject cannot be empty")
	}
	if plainText == "" {
		return "", fmt.Errorf("plainText cannot be empty")
	}

	msgJSON, err := json.Marshal(map[string]string{
		"htmlBody":  htmlBody,
		"plainText": plainText,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal SNS message: %w", err)
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Message:  aws.String(string(msgJSON)),
		Subject:  aws.String(subject),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"format": {
				DataType:    aws.String("String"),
				StringValue: aws.String("html"),
			},
		},
	}

	output, err := s.client.Publish(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to publish SNS message: %w", err)
	}

	return aws.ToString(output.MessageId), nil
}
