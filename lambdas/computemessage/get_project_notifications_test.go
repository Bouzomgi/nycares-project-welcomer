package main

import (
	"context"
	"errors"
	"testing"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// mockDynamoClient implements just the GetItem method
type mockDynamoClient struct {
	resp map[string]types.AttributeValue
	err  error
}

func (m *mockDynamoClient) GetItem(ctx context.Context, input *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &dynamodb.GetItemOutput{Item: m.resp}, nil
}

func TestGetProjectNotification_RowExists(t *testing.T) {
	mockResp := map[string]types.AttributeValue{
		"ProjectName":    &types.AttributeValueMemberS{Value: "ProjectA"},
		"ProjectDate":    &types.AttributeValueMemberS{Value: "2025-11-11"},
		"HasSentWelcome": &types.AttributeValueMemberBOOL{Value: true},
	}

	db := &mockDynamoClient{resp: mockResp}

	project := models.Project{Name: "ProjectA", Date: "2025-11-11"}
	pn, err := GetProjectNotification(db, "table", project)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pn == nil {
		t.Fatal("expected ProjectNotification, got nil")
	}
	if !pn.HasSentWelcome {
		t.Errorf("expected HasSentWelcome=true, got false")
	}
}

func TestGetProjectNotification_RowDoesNotExist(t *testing.T) {
	db := &mockDynamoClient{resp: map[string]types.AttributeValue{}}

	project := models.Project{Name: "ProjectB", Date: "2025-11-12"}
	pn, err := GetProjectNotification(db, "table", project)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if pn != nil {
		t.Errorf("expected nil ProjectNotification, got %+v", pn)
	}
}

func TestGetProjectNotification_DynamoError(t *testing.T) {
	db := &mockDynamoClient{err: errors.New("Dynamo failure")}

	project := models.Project{Name: "ProjectC", Date: "2025-11-13"}
	_, err := GetProjectNotification(db, "table", project)
	if err == nil || err.Error() != "dynamo get item failed: Dynamo failure" {
		t.Errorf("expected Dynamo failure error, got %v", err)
	}
}

func TestGetProjectNotification_EmptyTableName(t *testing.T) {
	db := &mockDynamoClient{}

	project := models.Project{Name: "ProjectD", Date: "2025-11-14"}
	_, err := GetProjectNotification(db, "", project)
	if err == nil || err.Error() != "dynamo table name is required" {
		t.Errorf("expected 'dynamo table name is required' error, got %v", err)
	}
}
