package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// publishMessage publishes a message to the specified SNS topic
func publishMessage(ctx context.Context, client *sns.Client, topicARN, message, subject string) error {
	// Validate input
	if topicARN == "" {
		return fmt.Errorf("topicARN cannot be empty")
	}
	if message == "" {
		return fmt.Errorf("message cannot be empty")
	}
	if subject == "" {
		return fmt.Errorf("subject cannot be empty")
	}

	input := &sns.PublishInput{
		TopicArn: aws.String(topicARN),
		Message:  aws.String(message),
		Subject:  aws.String(subject),
	}

	output, err := client.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to publish SNS message: %w", err)
	}

	log.Printf("Message published! MessageID: %s", *output.MessageId)
	return nil
}
