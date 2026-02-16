package snsservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type NotificationService interface {
	PublishNotification(ctx context.Context, message, subject string) (string, error)
}

// PublishNotification publishes a message to the configured SNS topic
func (s *SNSService) PublishNotification(ctx context.Context, message, subject string) (string, error) {
	// Validate input
	if s.topicARN == "" {
		return "", fmt.Errorf("topicARN cannot be empty")
	}
	if subject == "" {
		return "", fmt.Errorf("subject cannot be empty")
	}
	if message == "" {
		return "", fmt.Errorf("message cannot be empty")
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(s.topicARN),
		Message:  aws.String(message),
		Subject:  aws.String(subject),
	}

	output, err := s.client.Publish(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to publish SNS message: %w", err)
	}

	messageID := aws.ToString(output.MessageId)

	return messageID, nil
}
