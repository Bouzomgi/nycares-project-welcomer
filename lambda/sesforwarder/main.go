package main

import (
	"context"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	sesservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/ses"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

func buildHandler() (*SESForwarderHandler, error) {
	cfg, err := config.LoadConfig[Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	sesClient := sesv2.NewFromConfig(awsCfg)
	sesSvc := sesservice.NewSESService(sesClient, cfg.AWS.SES.Sender, cfg.AWS.SES.Recipient)

	return NewSESForwarderHandler(sesSvc), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		if err := handler.Handle(context.Background(), events.SNSEvent{}); err != nil {
			panic(err)
		}
		return
	}

	lambda.Start(handler.Handle)
}
