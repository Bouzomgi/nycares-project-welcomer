package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	cm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/computemessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func buildHandler() (*ComputeMessageHandler, error) {
	cfg, err := config.LoadConfig[cm.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	dynamoClient := dynamodb.NewFromConfig(awsCfg)
	dynamoSvc := dynamoservice.NewDynamoService(dynamoClient, cfg.AWS.Dynamo.TableName)

	var currentDate *time.Time
	if cfg.CurrentDate != "" {
		t, err := time.Parse("2006-01-02", cfg.CurrentDate)
		if err != nil {
			return nil, fmt.Errorf("invalid NYCARES_CURRENT_DATE %q: %w", cfg.CurrentDate, err)
		}
		currentDate = &t
	}

	usecase := cm.NewComputeMessageUseCase(dynamoSvc, currentDate)
	return NewComputeMessageHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		output, err := handler.Handle(context.Background(), models.ComputeMessageInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
