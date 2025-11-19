package main

import (
	"context"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/login"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
	"github.com/aws/aws-lambda-go/lambda"
)

func buildHandler() (*login.LoginHandler, error) {
	cfg, err := config.LoadConfig[login.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService()
	if err != nil {
		return nil, err
	}

	usecase := login.NewLoginUseCase(httpSvc)
	return login.NewLoginHandler(usecase, cfg), nil
}

func main() {
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		_, err := handler.Handle(context.Background())
		if err != nil {
			panic(err)
		}
		return
	}

	lambda.Start(handler.Handle)
}
