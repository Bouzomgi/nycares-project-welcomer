package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/login"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func buildHandler() (*LoginHandler, error) {
	cfg, err := config.LoadConfig[login.Config]()
	if err != nil {
		return nil, err
	}

	baseUrl := endpoints.BaseUrl
	if cfg.Api.BaseUrl != "" {
		baseUrl = cfg.Api.BaseUrl
	}

	httpSvc, err := httpservice.NewHttpService(baseUrl)
	if err != nil {
		return nil, err
	}

	usecase := login.NewLoginUseCase(httpSvc)
	return NewLoginHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		output, err := handler.Handle(context.Background())
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
