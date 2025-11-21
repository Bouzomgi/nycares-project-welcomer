package dynamoservice

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoService struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoService(client *dynamodb.Client, tableName string) *DynamoService {
	return &DynamoService{
		client:    client,
		tableName: tableName,
	}
}
