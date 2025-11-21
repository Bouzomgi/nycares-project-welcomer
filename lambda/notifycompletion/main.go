package main

import (
	"context"
	"os"

	nc "github.com/Bouzomgi/nycares-project-welcomer/internal/app/notifycompletion"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*NotifyCompletionHandler, error) {
	cfg, err := config.LoadConfig[nc.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsCfg)

	snsSvc := snsservice.NewSNSService(snsClient, cfg.AWS.SNS.TopicArn)

	usecase := nc.NewNotifyCompletionUseCase(snsSvc)
	return NewNotifyCompletionHandler(usecase, cfg), nil
}

func main() {
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		err := handler.Handle(context.Background(), models.NotifyCompletionInput{})
		if err != nil {
			panic(err)
		}
		return
	}

	lambda.Start(handler.Handle)
}
