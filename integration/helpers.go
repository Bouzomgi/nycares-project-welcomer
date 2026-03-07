//go:build integration

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	stateMachineARN  = "arn:aws:states:us-east-1:000000000000:stateMachine:project-notifier-workflow"
	dynamoTableName  = "Sent_Notifications"
	pollInterval     = 2 * time.Second
	executionTimeout = 120 * time.Second
)

type testClients struct {
	sfnClient     *sfn.Client
	dynamoClient  *dynamodb.Client
	mockServerURL string
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initTestClients() (*testClients, error) {
	endpoint := getEnvOrDefault("AWS_ENDPOINT_URL", "http://localhost:4566")
	mockServerURL := getEnvOrDefault("MOCKSERVER_URL", "http://localhost:3001")

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("test", "test", "")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
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
	}, nil
}

type projectInput struct {
	Name       string `json:"name"`
	Date       string `json:"date"`
	Id         string `json:"id"`
	CampaignId string `json:"campaignId"`
}

func (tc *testClients) setMockProjects(projects []projectInput) error {
	body, err := json.Marshal(map[string]interface{}{
		"projects": projects,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal projects: %w", err)
	}

	resp, err := http.Post(tc.mockServerURL+"/admin/set-projects", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to set mock projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("set-projects returned status %d", resp.StatusCode)
	}
	return nil
}

func (tc *testClients) startExecution() (string, error) {
	ctx := context.Background()
	result, err := tc.sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachineARN),
		Input:           aws.String("{}"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to start execution: %w", err)
	}
	return *result.ExecutionArn, nil
}

func (tc *testClients) waitForExecutionComplete(executionArn string) (sfntypes.ExecutionStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("execution timed out: %s", executionArn)
		default:
		}

		result, err := tc.sfnClient.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		})
		if err != nil {
			return "", fmt.Errorf("failed to describe execution: %w", err)
		}

		if result.Status != sfntypes.ExecutionStatusRunning {
			return result.Status, nil
		}

		time.Sleep(pollInterval)
	}
}

func (tc *testClients) pollForTaskToken(executionArn string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), executionTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timed out waiting for task token")
		default:
		}

		result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
			ExecutionArn: aws.String(executionArn),
			ReverseOrder: true,
		})
		if err != nil {
			return "", fmt.Errorf("failed to get execution history: %w", err)
		}

		for _, event := range result.Events {
			if event.TaskScheduledEventDetails != nil {
				params := aws.ToString(event.TaskScheduledEventDetails.Parameters)
				token := extractTaskToken(params)
				if token != "" {
					return token, nil
				}
			}
		}

		time.Sleep(pollInterval)
	}
}

func extractTaskToken(params string) string {
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(params), &parsed); err != nil {
		return ""
	}

	if payload, ok := parsed["Payload"].(map[string]interface{}); ok {
		if token, ok := payload["taskToken"].(string); ok {
			return token
		}
	}

	if token, ok := parsed["taskToken"].(string); ok {
		return token
	}

	return ""
}

func (tc *testClients) approveTask(taskToken string) error {
	ctx := context.Background()
	_, err := tc.sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
		TaskToken: aws.String(taskToken),
		Output:    aws.String(`{"approved": true}`),
	})
	if err != nil {
		return fmt.Errorf("failed to approve task: %w", err)
	}
	return nil
}

func (tc *testClients) rejectTask(taskToken string) error {
	ctx := context.Background()
	_, err := tc.sfnClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
		TaskToken: aws.String(taskToken),
		Error:     aws.String("rejected"),
		Cause:     aws.String("User rejected the approval request"),
	})
	if err != nil {
		return fmt.Errorf("failed to reject task: %w", err)
	}
	return nil
}

func (tc *testClients) getNotification(projectName, projectDate string) (map[string]types.AttributeValue, error) {
	ctx := context.Background()
	result, err := tc.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ProjectName": &types.AttributeValueMemberS{Value: projectName},
			"ProjectDate": &types.AttributeValueMemberS{Value: projectDate},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get DynamoDB item: %w", err)
	}
	return result.Item, nil
}

func (tc *testClients) deleteNotification(projectName, projectDate string) error {
	ctx := context.Background()
	_, err := tc.dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(dynamoTableName),
		Key: map[string]types.AttributeValue{
			"ProjectName": &types.AttributeValueMemberS{Value: projectName},
			"ProjectDate": &types.AttributeValueMemberS{Value: projectDate},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete DynamoDB item: %w", err)
	}
	return nil
}

func (tc *testClients) seedNotification(projectName, projectDate, projectId string, hasSentWelcome, hasSentReminder bool) error {
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
		return fmt.Errorf("failed to seed DynamoDB notification: %w", err)
	}
	return nil
}

func (tc *testClients) executionEndedWithSkip(executionArn string) (bool, error) {
	ctx := context.Background()
	result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get execution history: %w", err)
	}

	for _, event := range result.Events {
		if event.StateEnteredEventDetails != nil {
			name := aws.ToString(event.StateEnteredEventDetails.Name)
			if name == "EndProjectIteration" {
				return true, nil
			}
		}
	}
	return false, nil
}

func (tc *testClients) executionEnteredState(executionArn, stateName string) (bool, error) {
	ctx := context.Background()
	result, err := tc.sfnClient.GetExecutionHistory(ctx, &sfn.GetExecutionHistoryInput{
		ExecutionArn: aws.String(executionArn),
	})
	if err != nil {
		return false, fmt.Errorf("failed to get execution history: %w", err)
	}

	for _, event := range result.Events {
		if event.StateEnteredEventDetails != nil {
			if aws.ToString(event.StateEnteredEventDetails.Name) == stateName {
				return true, nil
			}
		}
	}
	return false, nil
}

func getBoolAttr(item map[string]types.AttributeValue, key string) (bool, error) {
	av, ok := item[key]
	if !ok {
		return false, fmt.Errorf("missing attribute %q", key)
	}
	boolAttr, ok := av.(*types.AttributeValueMemberBOOL)
	if !ok {
		return false, fmt.Errorf("attribute %q is not BOOL", key)
	}
	return boolAttr.Value, nil
}

func currentDateStr() string {
	if d := os.Getenv("NYCARES_CURRENT_DATE"); d != "" {
		return d
	}
	return time.Now().UTC().Format("2006-01-02")
}

func dateOffset(baseDate string, days int) string {
	t, err := time.Parse("2006-01-02", baseDate)
	if err != nil {
		panic(fmt.Sprintf("invalid base date: %s", baseDate))
	}
	return t.AddDate(0, 0, days).Format("2006-01-02")
}
