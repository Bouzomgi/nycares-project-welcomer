package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ac "github.com/Bouzomgi/nycares-project-welcomer/internal/app/approvalcallback"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*ApprovalCallbackHandler, error) {
	cfg, err := config.LoadConfig[ac.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	sfnClient := sfn.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)
	snsSvc := snsservice.NewSNSService(snsClient, cfg.AWS.SNS.TopicArn)

	usecase := ac.NewApprovalCallbackUseCase(sfnClient)
	return NewApprovalCallbackHandler(usecase, cfg, snsSvc), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		output, err := handler.Handle(context.Background(), events.APIGatewayProxyRequest{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
