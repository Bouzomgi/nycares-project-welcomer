package snsservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type NotificationService interface {
	PublishNotification(ctx context.Context, message, subject string) (string, error)
	PublishHTMLEmailNotification(ctx context.Context, plainText, htmlBody, subject string) (string, error)
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

// PublishHTMLEmailNotification publishes a JSON-structured message to the SNS topic,
// sending HTML to email subscribers and plain text to all other protocols.
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
		"default": plainText,
		"email":   htmlBody,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal SNS message: %w", err)
	}

	input := &sns.PublishInput{
		TopicArn:         aws.String(s.topicARN),
		Message:          aws.String(string(msgJSON)),
		Subject:          aws.String(subject),
		MessageStructure: aws.String("json"),
	}

	output, err := s.client.Publish(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to publish SNS message: %w", err)
	}

	return aws.ToString(output.MessageId), nil
}
