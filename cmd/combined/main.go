package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/fetchprojects"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/login"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
)

func buildLoginHandler() (*login.LoginHandler, error) {
	cfg, err := config.LoadConfig[login.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := login.NewLoginUseCase(httpSvc)
	return login.NewLoginHandler(usecase, cfg), nil
}

func buildFetchProjectsHandler() (*fetchprojects.FetchProjectsHandler, error) {
	cfg, err := config.LoadConfig[fetchprojects.Config]()
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	usecase := fetchprojects.NewFetchProjectsUseCase(httpSvc)
	return fetchprojects.NewFetchProjectsHandler(usecase, cfg), nil
}

////////////////

func main() {

	loginHandler, err := buildLoginHandler()
	if err != nil {
		panic(err)
	}

	loginOut, err := loginHandler.Handle(context.Background())
	if err != nil {
		panic(err)
	}

	data, _ := json.MarshalIndent(loginOut, "", "  ")
	fmt.Println(string(data))
	/////

	fetchProjectHandler, err := buildFetchProjectsHandler()
	if err != nil {
		panic(err)
	}

	output, err := fetchProjectHandler.Handle(context.Background(), loginOut)

	if err != nil {
		panic(err)
	}

	data, _ = json.MarshalIndent(output, "", "  ")
	fmt.Println(string(data))
	return
}
