package dynamoservice

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type ProjectNotification struct {
	ProjectName      string `json:"projectName"`
	ProjectDate      string `json:"projectDate"`
	HasSentWelcome   bool   `json:"hasSentWelcome"`
	HasSentReminder  bool   `json:"hasSentReminder"`
	ShouldStopNotify bool   `json:"shouldStopNotify"`
	LastUpdated      string `json:"lastUpdated"`
}

type StoredNotificationService interface {
	GetProjectNotification(ctx context.Context, project domain.Project) (*domain.ProjectNotification, error)
}

func (s *DynamoService) GetProjectNotification(ctx context.Context, project domain.Project) (*domain.ProjectNotification, error) {
	if s.tableName == "" {
		return nil, fmt.Errorf("dynamo table name is required")
	}

	key, err := attributevalue.MarshalMap(map[string]string{
		"ProjectName": project.Name,
		"ProjectDate": utils.DateToString(project.Date),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	resp, err := s.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})

	if err != nil {
		return nil, fmt.Errorf("dynamo get item failed: %w", err)
	}

	if len(resp.Item) == 0 {
		return nil, nil // row does not exist
	}

	var pn ProjectNotification
	if err := attributevalue.UnmarshalMap(resp.Item, &pn); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	domainNotification := &domain.ProjectNotification{
		ProjectName:      pn.ProjectName,
		ProjectDate:      pn.ProjectDate,
		HasSentWelcome:   pn.HasSentWelcome,
		HasSentReminder:  pn.HasSentReminder,
		ShouldStopNotify: pn.ShouldStopNotify,
	}

	return domainNotification, nil
}
