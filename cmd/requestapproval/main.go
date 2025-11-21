package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/requestapproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*requestapproval.RequestApprovalHandler, error) {
	cfg, err := config.LoadConfig[requestapproval.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsCfg)

	snsSvc := snsservice.NewSNSSerice(snsClient, cfg.AWS.SNS.TopicArn)

	usecase := requestapproval.NewRequestApprovalUseCase(snsSvc)
	return requestapproval.NewRequestApprovalHandler(usecase, cfg), nil
}

func main() {
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
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
