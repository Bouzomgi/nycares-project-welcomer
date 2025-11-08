package main

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// GetProjectNotification fetches a row by project name + date
func GetProjectNotification(dbClient *dynamodb.Client, dynamoTableName string, project models.Project) (*models.ProjectNotification, error) {
	if dynamoTableName == "" {
		log.Errorf("Dynamo table name is empty for project %s", project.Name)
		return nil, fmt.Errorf("dynamo table name is required")
	}

	key, err := attributevalue.MarshalMap(map[string]string{
		"ProjectName": project.Name,
		"ProjectDate": project.Date,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal key: %w", err)
	}

	resp, err := dbClient.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(dynamoTableName),
		Key:       key,
	})
	if err != nil {
		return nil, fmt.Errorf("dynamo get item failed: %w", err)
	}

	if len(resp.Item) == 0 {
		return nil, nil // row does not exist
	}

	var pn models.ProjectNotification
	if err := attributevalue.UnmarshalMap(resp.Item, &pn); err != nil {
		return nil, fmt.Errorf("failed to unmarshal item: %w", err)
	}

	return &pn, nil
}
