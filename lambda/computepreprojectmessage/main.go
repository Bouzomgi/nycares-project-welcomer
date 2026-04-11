package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	cm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/computepreprojectmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
)

func buildHandler() (*ComputeMessageHandler, error) {
	cfg, err := config.LoadConfig[cm.Config]()
	if err != nil {
		return nil, err
	}

	usecase := cm.NewComputeMessageUseCase()
	return NewComputeMessageHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		var input models.ComputeMessageInput
		if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
			panic(fmt.Errorf("failed to decode input: %w", err))
		}

		output, err := handler.Handle(context.Background(), input)
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
