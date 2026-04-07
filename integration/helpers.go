//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"fmt"
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

const pollInterval = 2 * time.Second

var (
	stateMachineARN = getEnvOrDefault("STATE_MACHINE_ARN", "arn:aws:states:us-east-1:000000000000:stateMachine:project-notifier-workflow")
	dynamoTableName = getEnvOrDefault("DYNAMO_TABLE_NAME", "nycares-project-welcomer-notifications")
)

var executionTimeout = func() time.Duration {
	if v := os.Getenv("INTEGRATION_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return 20 * time.Second
}()

type testClients struct {
	sfnClient    *sfn.Client
	dynamoClient *dynamodb.Client
}

func getEnvOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func initTestClients() (*testClients, error) {
	awsEndpoint := os.Getenv("AWS_ENDPOINT_URL") // empty = use real AWS
	ctx := context.Background()

	var cfgOpts []func(*config.LoadOptions) error
	cfgOpts = append(cfgOpts, config.WithRegion("us-east-1"))
	if awsEndpoint != "" {
		// LocalStack: use static test credentials
		cfgOpts = append(cfgOpts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		))
	}

	cfg, err := config.LoadDefaultConfig(ctx, cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	var sfnOpts []func(*sfn.Options)
	var dynamoOpts []func(*dynamodb.Options)
	if awsEndpoint != "" {
		sfnOpts = append(sfnOpts, func(o *sfn.Options) {
			o.BaseEndpoint = aws.String(awsEndpoint)
		})
		dynamoOpts = append(dynamoOpts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(awsEndpoint)
		})
	}

	sfnClient := sfn.NewFromConfig(cfg, sfnOpts...)
	dynamoClient := dynamodb.NewFromConfig(cfg, dynamoOpts...)

	return &testClients{
		sfnClient:    sfnClient,
		dynamoClient: dynamoClient,
	}, nil
}

type projectInput struct {
	Name   string `json:"name"`
	Date   string `json:"date"`
	Id     string `json:"id"`
	Status string `json:"status,omitempty"`
}

func (tc *testClients) startExecutionWithProjects(projects []projectInput) (string, error) {
	ctx := context.Background()
	input, err := json.Marshal(map[string]interface{}{"mockProjects": projects})
	if err != nil {
		return "", fmt.Errorf("failed to marshal projects: %w", err)
	}
	result, err := tc.sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachineARN),
		Input:           aws.String(string(input)),
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
