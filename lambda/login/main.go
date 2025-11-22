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

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := login.NewLoginUseCase(httpSvc)
	return NewLoginHandler(usecase, cfg), nil
}

func main() {
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
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
