//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	sfntypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

const (
	stateMachineARN = "arn:aws:states:us-east-1:000000000000:stateMachine:project-notifier-workflow"
	dynamoTableName = "Sent_Notifications"
	pollInterval    = 2 * time.Second
	executionTimeout = 120 * time.Second
)

type testClients struct {
	sfnClient    *sfn.Client
	dynamoClient *dynamodb.Client
	mockServerURL string
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newTestClients(t *testing.T) *testClients {
	t.Helper()

	endpoint := getEnvOrDefault("AWS_ENDPOINT_URL", "http://localhost:4566")
	mockServerURL := getEnvOrDefault("MOCKSERVER_URL", "http://localhost:3001")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	sfnClient := sfn.NewFromConfig(cfg, func(o *sfn.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	dynamoClient := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &testClients{
		sfnClient:     sfnClient,
		dynamoClient:  dynamoClient,
		mockServerURL: mockServerURL,
	}
}

type projectInput struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Id         string `json:"id"`
	CampaignId string `json:"campaignId"`
}

func (tc *testClients) setMockProjects(t *testing.T, projects []projectInput) {
	t.Helper()

	body, err := json.Marshal(map[string]interface{}{
		"projects": projects,
	})
	if err != nil {
		t.Fatalf("failed to marshal projects: %v", err)
	}

	resp, err := http.Post(tc.mockServerURL+"/admin/set-projects", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("failed to set mock projects: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("set-projects returned status %d", resp.StatusCode)
	}
}

func (tc *testClients) startExecution(t *testing.T) string {
	t.Helper()

	ctx := context.Background()
	result, err := tc.sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachineARN),
		Input:           aws.String("{}"),
	})
	if err != nil {
		t.Fatalf("failed to start execution: %v", err)
	}

	return *result.ExecutionArn
}

func (tc *testClients) waitForExecutionComplete(t *testing.T, executionArn string) sfntypes.ExecutionStatus {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("execution timed out: %s", executionArn)
		default:
		}

		result, err := tc.sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		})
		if err != nil {
			t.Fatalf("failed to describe execution: %v", err)
		}

		if result.Status != sfntypes.ExecutionStatusRunning {
			return result.Status
		}

		time.Sleep(pollInterval)
	}
}

func (tc *testClients) pollForTaskToken(t *testing.T, executionArn string) string {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			t.Fatalf("timed out waiting for task token")
		default:
		}

		result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(executionArn),
			ReverseOrder: true,
		})
		if err != nil {
			t.Fatalf("failed to get execution history: %v", err)
		}

		for _, event := range result.Events {
			if event.TaskScheduledEventDetails != nil {
				// The task token is in the Parameters field of TaskScheduled events
				// for waitForTaskToken integrations
				params := aws.ToString(event.TaskScheduledEventDetails.Parameters)
				token := extractTaskToken(params)
				if token != "" {
					return token
				}
			}
		}

		time.Sleep(pollInterval)
	}
}

func extractTaskToken(params string) string {
	// Parameters is a JSON string containing the payload with taskToken
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(params), &parsed); err != nil {
		return ""
	}

	// Look for taskToken in Payload
	if payload, ok := parsed["Payload"].(map[string]interface{}); ok {
		if token, ok := payload["taskToken"].(string); ok {
			return token
		}
	}

	// Also check top level
	if token, ok := parsed["taskToken"].(string); ok {
		return token
	}

	return ""
}

func (tc *testClients) approveTask(t *testing.T, taskToken string) {
	t.Helper()

	ctx := context.Background()
	_, err := tc.sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
		TaskToken: aws.String(taskToken),
		Output:    aws.String(`{"approved": true}`),
	})
	if err != nil {
		t.Fatalf("failed to approve task: %v", err)
	}
}

func (tc *testClients) rejectTask(t *testing.T, taskToken string) {
	t.Helper()

	ctx := context.Background()
	_, err := tc.sfnClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
		TaskToken: aws.String(taskToken),
		Error:     aws.String("rejected"),
		Cause:     aws.String("User rejected the approval request"),
	})
	if err != nil {
		t.Fatalf("failed to reject task: %v", err)
	}
}

func (tc *testClients) getNotification(t *testing.T, projectName, projectDate string) *map[string]types.AttributeValue {
	t.Helper()

	ctx := context.Background()
	result, err := tc.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ProjectName": &types.AttributeValueMemberS{Value: projectName},
			"ProjectDate": &types.AttributeValueMemberS{Value: projectDate},
		},
	})
	if err != nil {
		t.Fatalf("failed to get DynamoDB item: %v", err)
	}

	if result.Item == nil {
		return nil
	}
	return &result.Item
}

func (tc *testClients) deleteNotification(t *testing.T, projectName, projectDate string) {
	t.Helper()

	ctx := context.Background()
	_, err := tc.dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ProjectName": &types.AttributeValueMemberS{Value: projectName},
			"ProjectDate": &types.AttributeValueMemberS{Value: projectDate},
		},
	})
	if err != nil {
		t.Fatalf("failed to delete DynamoDB item: %v", err)
	}
}

func (tc *testClients) seedNotification(t *testing.T, projectName, projectDate, projectId string, hasSentWelcome, hasSentReminder bool) {
	t.Helper()

	ctx := context.Background()
	_, err := tc.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(dynamoTableName),
		Item: map[string]types.AttributeValue{
			"ProjectName":      &types.AttributeValueMemberS{Value: projectName},
			"ProjectDate":      &types.AttributeValueMemberS{Value: projectDate},
			"ProjectId":        &types.AttributeValueMemberS{Value: projectId},
			"HasSentWelcome":   &types.AttributeValueMemberBOOL{Value: hasSentWelcome},
			"HasSentReminder":  &types.AttributeValueMemberBOOL{Value: hasSentReminder},
			"ShouldStopNotify": &types.AttributeValueMemberBOOL{Value: false},
			"LastUpdated":      &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	})
	if err != nil {
		t.Fatalf("failed to seed DynamoDB notification: %v", err)
	}
}

func assertBoolAttr(t *testing.T, item map[string]types.AttributeValue, key string, expected bool) {
	t.Helper()
	av, ok := item[key]
	if !ok {
		t.Fatalf("missing attribute %q", key)
	}
	boolAttr, ok := av.(*types.AttributeValueMemberBOOL)
	if !ok {
		t.Fatalf("attribute %q is not BOOL", key)
	}
	if boolAttr.Value != expected {
		t.Errorf("attribute %q = %v, want %v", key, boolAttr.Value, expected)
	}
}

func (tc *testClients) getExecutionOutput(t *testing.T, executionArn string) string {
	t.Helper()

	ctx := context.Background()
	result, err := tc.sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		t.Fatalf("failed to describe execution: %v", err)
	}
	return aws.ToString(result.Output)
}

func (tc *testClients) executionEndedWithSkip(t *testing.T, executionArn string) bool {
	t.Helper()

	ctx := context.Background()
	result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		t.Fatalf("failed to get execution history: %v", err)
	}

	for _, event := range result.Events {
		if event.StateEnteredEventDetails != nil {
			name := aws.ToString(event.StateEnteredEventDetails.Name)
			if name == "EndProjectIteration" {
				return true
			}
		}
	}
	return false
}

func (tc *testClients) executionEnteredState(t *testing.T, executionArn, stateName string) bool {
	t.Helper()

	ctx := context.Background()
	result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		t.Fatalf("failed to get execution history: %v", err)
	}

	for _, event := range result.Events {
		if event.StateEnteredEventDetails != nil {
			if aws.ToString(event.StateEnteredEventDetails.Name) == stateName {
				return true
			}
		}
	}
	return false
}

func currentDateStr() string {
	if d := os.Getenv("NYCARES_CURRENT_DATE"); d != "" {
		return d
	}
	return time.Now().Format("2006-01-02")
}

func dateOffset(baseDate string, days int) string {
	t, err := time.Parse("2006-01-02", baseDate)
	if err != nil {
		panic(fmt.Sprintf("invalid base date: %s", baseDate))
	}
	return t.AddDate(0, 0, days).Format("2006-01-02")
}

func executionHasError(t *testing.T, tc *testClients, executionArn string, errorSubstr string) bool {
	t.Helper()

	ctx := context.Background()
	result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		t.Fatalf("failed to get execution history: %v", err)
	}

	for _, event := range result.Events {
		if event.LambdaFunctionFailedEventDetails != nil {
			cause := aws.ToString(event.LambdaFunctionFailedEventDetails.Cause)
			if strings.Contains(cause, errorSubstr) {
				return true
			}
		}
		if event.TaskFailedEventDetails != nil {
			cause := aws.ToString(event.TaskFailedEventDetails.Cause)
			if strings.Contains(cause, errorSubstr) {
				return true
			}
		}
	}
	return false
}
