package dynamoservice

import (
	"context"
	"fmt"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/dto"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type StoredNotificationService interface {
	GetProjectNotification(ctx context.Context, project domain.Project) (*domain.ProjectNotification, error)
	UpsertProjectNotification(ctx context.Context, projectNotification domain.ProjectNotification) (*domain.ProjectNotification, error)
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

	resp, err := s.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key:       key,
	})

	if err != nil {
		return nil, fmt.Errorf("dynamo get item failed: %w", err)
	}

	if len(resp.Item) == 0 {
		return nil, nil // row does not exist
	}

	var pn dto.ProjectNotification
	if err := attributevalue.UnmarshalMap(resp.Item, &pn); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	projectDate, err := utils.StringToDate(pn.ProjectDate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse project date: %w", err)
	}

	domainNotification := &domain.ProjectNotification{
		Name:             pn.ProjectName,
		Date:             projectDate,
		Id:               pn.ProjectId,
		HasSentWelcome:   pn.HasSentWelcome,
		HasSentReminder:  pn.HasSentReminder,
		ShouldStopNotify: pn.ShouldStopNotify,
	}

	return domainNotification, nil
}

func (s *DynamoService) UpsertProjectNotification(
	ctx context.Context,
	projectNotification domain.ProjectNotification,
) (*domain.ProjectNotification, error) {

	if s.tableName == "" {
		return nil, fmt.Errorf("dynamo table name is required")
	}

	now := time.Now().UTC().Format(time.RFC3339)

	item := dto.ProjectNotification{
		ProjectName:      projectNotification.Name,
		ProjectDate:      utils.DateToString(projectNotification.Date),
		ProjectId:        projectNotification.Id,
		HasSentWelcome:   projectNotification.HasSentWelcome,
		HasSentReminder:  projectNotification.HasSentReminder,
		ShouldStopNotify: projectNotification.ShouldStopNotify,
		LastUpdated:      now,
	}

	av, err := attributevalue.MarshalMap(item)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal item: %w", err)
	}

	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.tableName),
		Item:      av,
	})
	if err != nil {
		return nil, fmt.Errorf("dynamo put item failed: %w", err)
	}

	result := &domain.ProjectNotification{
		Name:             item.ProjectName,
		Date:             projectNotification.Date,
		Id:               item.ProjectId,
		HasSentWelcome:   item.HasSentWelcome,
		HasSentReminder:  item.HasSentReminder,
		ShouldStopNotify: item.ShouldStopNotify,
	}

	return result, nil
}
