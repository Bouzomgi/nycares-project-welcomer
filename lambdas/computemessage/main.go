package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"

	"github.com/aws/aws-lambda-go/lambda"
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

	if err := ValidateConfig(appCfg, isLocal); err != nil {
		log.Fatalf("Invalid config: %v", err)
	}

	dbClient, err = InitAWSClients(appCfg, isLocal)
	if err != nil {
		log.Fatalf("Failed to initialize DynamoDB client: %v", err)
	}

	log.Infof("Initialized DynamoDB client for region %s, table %s", appCfg.AWS.Dynamo.Region, appCfg.AWS.Dynamo.TableName)
}

// --- Lambda Handler ---

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {
	projectDate, err := parseProjectDate(input.Project.Date)
	if err != nil {
		return LambdaOutput{}, fmt.Errorf("invalid project date: %w", err)
	}
	log.Infof("Parsed project date for %s: %s", input.Project.Name, projectDate.Format("2006-01-02"))

	err = CheckProjectDate(input.Project)
	if err != nil {
		return LambdaOutput{}, err
	}

	log.Infof("Fetching project notification from DynamoDB for %s on %s", input.Project.Name, input.Project.Date)
	existingRow, err := GetProjectNotification(dbClient, appCfg.AWS.Dynamo.TableName, input.Project)
	if err != nil {
		log.Errorf("error fetching project notification: %v", err)
		return LambdaOutput{}, err
	}

	row, sendable, err := HandleProjectNotification(appCfg.AWS.S3.BucketName, input.Project, existingRow)
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

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	lambda.Start(handler)
}
