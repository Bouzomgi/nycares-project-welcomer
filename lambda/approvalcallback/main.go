package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ac "github.com/Bouzomgi/nycares-project-welcomer/internal/app/approvalcallback"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
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

	usecase := ac.NewApprovalCallbackUseCase(sfnClient)
	return NewApprovalCallbackHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
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
