package dynamoservice

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBClient interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type DynamoService struct {
	client    DynamoDBClient
	tableName string
}

func NewDynamoService(cfg aws.Config, tableName string) *DynamoService {
	return &DynamoService{
		client:    dynamodb.NewFromConfig(cfg),
		tableName: tableName,
	}
}
