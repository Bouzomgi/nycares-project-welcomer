package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// --- Lambda Input/Output ---

type LambdaInput struct {
	Auth    models.Auth    `json:"auth"`
	Project models.Project `json:"project"`
}

type LambdaOutput struct {
	Auth                models.Auth                 `json:"auth"`
	Project             models.Project              `json:"project"`
	SendableMessage     models.SendableMessage      `json:"messageToSend"`
	ProjectNotification *models.ProjectNotification `json:"savedProjectNotification,omitempty"`
}

// --- Global Variables ---

var appCfg *Config
var dbClient *dynamodb.Client

// --- Initialization ---

func init() {
	isLocal := os.Getenv("_LAMBDA_SERVER_PORT") == ""

	var err error
	appCfg, err = confighelper.LoadConfig[Config]()
	if err != nil {
		log.Fatalf("Failed to load config.yaml: %v", err)
	}

	if err := validateConfig(appCfg, isLocal); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	dbClient, err = initAWSClients(appCfg, isLocal)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}

	log.Infof("Initialized DynamoDB client for region %s, table %s", appCfg.AWS.Dynamo.Region, appCfg.AWS.Dynamo.TableName)
}

// --- AWS Client Setup ---

func initAWSClients(cfg *Config, isLocal bool) (*dynamodb.Client, error) {
	var awsCfg aws.Config
	var err error

	if isLocal {
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.AWS.Dynamo.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cfg.AWS.Credentials.AccessKeyID,
					cfg.AWS.Credentials.SecretAccessKey,
					"",
				),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS SDK config for local: %w", err)
		}

		// DynamoDB local endpoint
		endpoint := cfg.AWS.Dynamo.Endpoint
		if endpoint == "" {
			endpoint = "http://localhost:8000"
		}

		return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &endpoint
		}), nil
	}

	// Lambda: use IAM role automatically
	awsCfg, err = config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Credentials.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for Lambda: %w", err)
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}

// --- Config Validation ---

func validateConfig(cfg *Config, isLocal bool) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if cfg.AWS.Credentials.Region == "" {
		return fmt.Errorf("AWS region is missing")
	}

	if isLocal {
		if cfg.AWS.Credentials.AccessKeyID == "" || cfg.AWS.Credentials.SecretAccessKey == "" {
			return fmt.Errorf("local AWS credentials missing")
		}
		if cfg.AWS.Dynamo.Endpoint == "" {
			log.Warn("DynamoDB local endpoint not set, using default http://localhost:8000")
		}
		if cfg.AWS.S3.Endpoint == "" {
			log.Warn("S3 local endpoint not set")
		}
	}

	return nil
}

// --- Lambda Handler ---

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {
	projectDate, err := parseProjectDate(input.Project.Date)
	if err != nil {
		return LambdaOutput{}, err
	}

	daysUntilProject := time.Until(projectDate).Hours() / 24
	if daysUntilProject > 7 {
		log.Infof("Project %s is more than a week away, stopping execution", input.Project.Name)
		return LambdaOutput{}, fmt.Errorf("ProjectTooFar")
	}

	if daysUntilProject <= 0 {
		log.Infof("Project %s has passed, stopping execution", input.Project.Name)
		return LambdaOutput{}, fmt.Errorf("ProjectPassed")
	}

	existingRow, err := GetProjectNotification(dbClient, appCfg.AWS.Dynamo.TableName, input.Project)
	if err != nil {
		log.Errorf("error fetching project notification: %v", err)
		return LambdaOutput{}, err
	}

	row, sendable, err := handleProjectNotification(appCfg.AWS.S3.BucketName, input.Project, existingRow)
	if err != nil {
		return LambdaOutput{}, err
	}

	output := LambdaOutput{
		Auth:                input.Auth,
		Project:             input.Project,
		SendableMessage:     sendable,
		ProjectNotification: row,
	}

	// ✅ Print the final output
	log.Infof("Lambda output: %+v", output)

	return output, nil
}

// --- Helpers ---

func parseProjectDate(dateStr string) (time.Time, error) {
	projectDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Errorf("invalid project date format: %v", err)
		return time.Time{}, err
	}
	return projectDate, nil
}

func handleProjectNotification(s3BucketName string, project models.Project, existingNotification *models.ProjectNotification) (*models.ProjectNotification, models.SendableMessage, error) {
	if s3BucketName == "" {
		log.Errorf("S3 bucket name is empty for project %s", project.Name)
		return nil, models.SendableMessage{}, fmt.Errorf("S3 bucket name is required")
	}

	// No existing notification yet
	if existingNotification == nil {
		notification := createNewNotification(s3BucketName, project, models.Welcome)
		return nil, notification, nil
	}

	if existingNotification.ShouldStopNotify {
		log.Infof("Notifications are disabled for project %s", project.Name)
		return nil, models.SendableMessage{}, fmt.Errorf("ShouldStopNotify is true")
	}

	if !existingNotification.HasSentWelcome {
		notification := createNewNotification(s3BucketName, project, models.Welcome)
		return existingNotification, notification, nil
	}

	if !existingNotification.HasSentReminder {
		notification := createNewNotification(s3BucketName, project, models.Reminder)
		return existingNotification, notification, nil
	}

	log.Infof("All notifications already sent for project %s", project.Name)
	return nil, models.SendableMessage{}, fmt.Errorf("AllNotificationsAlreadySent")
}

func createNewNotification(s3BucketName string, project models.Project, messageType models.MessageType) models.SendableMessage {
	messageTypeStr := strings.ToLower(messageType.String())
	s3Destination := fmt.Sprintf("s3://%s/%s/%s.md", s3BucketName, ToCamelCase(project.Name), messageTypeStr)

	return models.SendableMessage{
		Type:        messageTypeStr,
		TemplateRef: s3Destination,
	}
}

// --- Main ---

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		testInput := LambdaInput{
			Auth: models.Auth{},
			Project: models.Project{
				Name: "ProjectA",
				Date: "2025-11-09",
			},
		}
		handler(context.Background(), testInput)
		return
	}

	lambda.Start(handler)
}
