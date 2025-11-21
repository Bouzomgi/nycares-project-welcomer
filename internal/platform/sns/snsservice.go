package snsservice

import (
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

// SNSService handles SNS notification operations
type SNSService struct {
	client   *sns.Client
	topicARN string
}

// NewSNSService creates a new SNS notification service
func NewSNSService(client *sns.Client, topicARN string) *SNSService {
	return &SNSService{
		client:   client,
		topicARN: topicARN,
	}
}
