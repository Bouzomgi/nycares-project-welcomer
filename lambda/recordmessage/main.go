package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	rm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/recordmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	dynamoservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/dynamo/service"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func buildHandler() (*RecordMessageHandler, error) {
	cfg, err := config.LoadConfig[rm.Config]()
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

	usecase := rm.NewRecordMessageUseCase(dynamoSvc)
	return NewRecordMessageHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		output, err := handler.Handle(context.Background(), models.RecordMessageInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
