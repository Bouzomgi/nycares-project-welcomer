package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ra "github.com/Bouzomgi/nycares-project-welcomer/internal/app/requestapproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*RequestApprovalHandler, error) {
	cfg, err := config.LoadConfig[ra.Config]()
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

	usecase := ra.NewRequestApprovalUseCase(snsSvc, cfg.AWS.SF.ApprovalSecret)
	return NewRequestApprovalHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		output, err := handler.Handle(context.Background(), models.RequestApprovalInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
