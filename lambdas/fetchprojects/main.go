package main

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/httphelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

// LambdaInput is what this Lambda receives from the previous step.
type LambdaInput struct {
	Auth models.Auth `json:"auth"`
}

// LambdaOutput defines what this Lambda returns.
type LambdaOutput struct {
	Auth     models.Auth `json:"auth"`
	Projects []Project   `json:"projects"`
}

type Project struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

// handler runs on Lambda invocation.
func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {

	cfg, err := confighelper.LoadConfig[Config]()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	if err := httphelper.SetCookiesOnClient(client, endpoints.BaseUrl, input.Auth); err != nil {
		log.Errorf("Failed to set cookies: %v", err)
		return LambdaOutput{}, err
	}

	log.Info("Fetching schedule")
	schedule, err := GetSchedule(client, endpoints.BaseUrl, cfg.Account.InternalId)
	if err != nil {
		log.Errorf("GetSchedule failed: %v", err)
		return LambdaOutput{}, err
	}

	log.Infof("Fetched schedule, has %d events", len(schedule))

	var projects []Project
	for _, s := range schedule {
		projects = append(projects, reduceToProject(s))
	}

	return LambdaOutput{
		Auth:     input.Auth,
		Projects: projects,
	}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background(), LambdaInput{})
		return
	}

	// Lambda runtime
	lambda.Start(handler)
}
