package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	fp "github.com/Bouzomgi/nycares-project-welcomer/internal/app/fetchprojects"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
	"github.com/aws/aws-lambda-go/lambda"
)

func buildHandler() (*FetchProjectsHandler, error) {
	cfg, err := config.LoadConfig[fp.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := fp.NewFetchProjectsUseCase(httpSvc)
	return NewFetchProjectsHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		output, err := handler.Handle(context.Background(), models.FetchProjectsInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
