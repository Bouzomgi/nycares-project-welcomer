package main

import (
	"context"
	"log"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
)

type LambdaInput struct {
	Auth                models.Auth                 `json:"auth"`
	Project             models.Project              `json:"project"`
	SendableMessage     models.SendableMessage      `json:"messageToSend"`
	ProjectNotification *models.ProjectNotification `json:"savedProjectNotification,omitempty"`
}

type LambdaOutput struct {
	Auth                models.Auth                 `json:"auth"`
	Project             models.Project              `json:"project"`
	SendableMessage     models.SendableMessage      `json:"messageToSend"`
	ProjectNotification *models.ProjectNotification `json:"savedProjectNotification,omitempty"`
}

var (
	appConfig *Config
)

func init() {
	var err error
	appConfig, err = confighelper.LoadConfig[Config]()
	if err != nil {
		log.Fatalf("Failed to load config.yaml: %v", err)
	}
}

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {

	// template := GetTemplate(appConfig.AWS.S3.BucketName, input.SendableMessage.TemplateRef)

	// SendMessage()

	// PinMessage()

	return LambdaOutput(input), nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	lambda.Start(handler)
}
