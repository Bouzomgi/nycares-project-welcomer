package computemessage

import (
	"context"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/computemessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/dynamoservice"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/s3service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"

	log "github.com/sirupsen/logrus"
)

type LambdaInput struct {
	Auth    models.Auth    `json:"auth"`
	Project models.Project `json:"project"`
}

type LambdaOutput struct {
	Auth                 models.Auth                 `json:"auth"`
	Project              models.Project              `json:"project"`
	ExistingNotification *models.ProjectNotification `json:"existingNotification,omitempty"`
	SendableMessage      models.SendableMessage      `json:"sendableMessage"`
}

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig(os.Getenv("NYCARES_ENV"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {
	awsCfg := *aws.NewConfig()

	dynamoService := dynamoservice.NewDynamoService(awsCfg, cfg.AWS.Dynamo.TableName)
	existingNotification, err := dynamoService.GetProjectNotification(ctx, input.Project)
	if err != nil {
		log.Errorf("Failed to get project notification: %v", err)
		return LambdaOutput{}, err
	}

	s3Service := s3service.NewS3Service(awsCfg, cfg.AWS.S3.BucketName)

	sendableMessage, err := computemessage.ComputeProjectMessage(
		s3Service,
		input.Project,
		existingNotification,
	)
	if err != nil {
		log.Errorf("Failed to compute project message: %v", err)
		return LambdaOutput{}, err
	}

	return LambdaOutput{
		Auth:                 input.Auth,
		Project:              input.Project,
		ExistingNotification: existingNotification,
		SendableMessage:      sendableMessage,
	}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	lambda.Start(handler)
}
