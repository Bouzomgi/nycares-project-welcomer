package fetchprojects

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	requestapproval "github.com/Bouzomgi/nycares-project-welcomer/internal/requestApproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/snsservice"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type LambdaInput struct {
	Auth                 models.Auth                 `json:"auth"`
	Project              models.Project              `json:"project"`
	SendableMessage      models.SendableMessage      `json:"messageToSend"`
	ExistingNotification *models.ProjectNotification `json:"existingNotification,omitempty"`
	TaskToken            string                      `json:"taskToken"`
}

type LambdaOutput struct {
	Auth                 models.Auth                 `json:"auth"`
	Project              models.Project              `json:"project"`
	SendableMessage      models.SendableMessage      `json:"messageToSend"`
	ExistingNotification *models.ProjectNotification `json:"existingNotification,omitempty"`
	CallbackLink         string                      `json:"callbackLink"`
}

var (
	cfg       *config.Config
	snsClient *sns.Client
)

func init() {
	var err error
	cfg, err = config.LoadConfig(os.Getenv("NYCARES_ENV"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {
	notificationService := snsservice.NewSNSSerice(snsClient, cfg.AWS.SNS.TopicArn)

	callbackURL := fmt.Sprintf(
		"%s?taskToken=%s",
		cfg.AWS.SF.CallbackEndpoint,
		input.TaskToken,
	)

	notificationMessage := requestapproval.CreateNotificationMessage(
		input.SendableMessage.Type,
		input.Project.Name,
		input.Project.Date,
		input.SendableMessage.TemplateRef,
	)

	notificationSubject := requestapproval.CreateNotificationSubject(
		input.SendableMessage.Type,
		input.Project.Name,
		input.Project.Date,
	)

	_, err := notificationService.PublishMessage(ctx, notificationMessage, notificationSubject)
	if err != nil {
		return LambdaOutput{}, fmt.Errorf("failed to publish notification: %w", err)
	}

	return LambdaOutput{
		Auth:                 input.Auth,
		Project:              input.Project,
		SendableMessage:      input.SendableMessage,
		ExistingNotification: input.ExistingNotification,
		CallbackLink:         callbackURL,
	}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	lambda.Start(handler)
}
