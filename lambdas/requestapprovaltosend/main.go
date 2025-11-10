package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type LambdaInput struct {
	Auth                models.Auth                 `json:"auth"`
	Project             models.Project              `json:"project"`
	SendableMessage     models.SendableMessage      `json:"messageToSend"`
	ProjectNotification *models.ProjectNotification `json:"savedProjectNotification,omitempty"`
	TaskToken           string                      `json:"taskToken"`
}

type LambdaOutput struct {
	Auth                models.Auth                 `json:"auth"`
	Project             models.Project              `json:"project"`
	SendableMessage     models.SendableMessage      `json:"messageToSend"`
	ProjectNotification *models.ProjectNotification `json:"savedProjectNotification,omitempty"`
	Link                string                      `json:"link"`
}

var (
	appConfig *Config
	snsClient *sns.Client
)

func init() {
	var err error
	// Load your custom config.yaml
	appConfig, err = confighelper.LoadConfig[Config]()
	if err != nil {
		log.Fatalf("Failed to load config.yaml: %v", err)
	}

	// Load AWS SDK config
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS SDK config: %v", err)
	}

	// Create SNS client
	snsClient = sns.NewFromConfig(awsCfg)
}

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {
	// Log or send the approval request externally
	fmt.Printf("Requesting approval for project %s with template %s\n", input.Project.Name, input.SendableMessage.TemplateRef)

	// Return a callback link
	callbackURL := fmt.Sprintf(
		"%s?taskToken=%s",
		appConfig.CallbackEndpoint,
		input.TaskToken,
	)

	notificationMessage := fmt.Sprintf(
		"Requesting approval to publish %s message to project %s on %s with template %s",
		input.SendableMessage.Type,
		input.Project.Name,
		input.Project.Date,
		input.SendableMessage.TemplateRef,
	)

	notificationSubject := fmt.Sprintf(
		"Send %s to %s on %s?",
		input.SendableMessage.Type,
		input.Project.Name,
		input.Project.Date,
	)

	publishMessage(ctx, snsClient, appConfig.AWS.SNS.TopicARN, notificationMessage, notificationSubject)

	output := LambdaOutput{
		Auth:                input.Auth,
		Project:             input.Project,
		SendableMessage:     input.SendableMessage,
		ProjectNotification: input.ProjectNotification,
		Link:                callbackURL,
	}

	return output, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	lambda.Start(handler)
}
