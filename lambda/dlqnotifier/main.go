package main

import (
	"context"
	"os"

	dlq "github.com/Bouzomgi/nycares-project-welcomer/internal/app/dlqnotifier"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*DLQNotifierHandler, error) {
	cfg, err := config.LoadConfig[dlq.Config]()
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

	usecase := dlq.NewDLQNotifierUseCase(snsSvc)
	return NewDLQNotifierHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		err := handler.Handle(context.Background(), models.DLQNotifierInput{})
		if err != nil {
			panic(err)
		}
		return
	}

	lambda.Start(handler.Handle)
}
