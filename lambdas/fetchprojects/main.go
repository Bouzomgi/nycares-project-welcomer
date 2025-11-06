package main

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

// LambdaOutput defines what the Lambda returns.
type LambdaOutput struct {
	Cookies  map[string]string `json:"cookies"`
	Schedule []models.Project  `json:"schedule"`
}

// handler runs on Lambda invocation.
func handler(ctx context.Context) (LambdaOutput, error) {

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	creds := Credentials{
		Username: cfg.Account.Username,
		Password: cfg.Account.Password,
	}

	log.Info("Attempting login")
	err = Login(client, PostLoginUrl, creds)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		return LambdaOutput{}, err
	}

	log.Info("Fetching schedule")
	schedule, err := GetSchedule(client, cfg.Account.InternalId)
	if err != nil {
		log.Errorf("GetSchedule failed: %v", err)
		return LambdaOutput{}, err
	}

	log.Infof("Fetched schedule, has %d events", len(schedule))

	return LambdaOutput{
		Schedule: schedule,
	}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		handler(context.Background())
		return
	}

	// Lambda runtime
	lambda.Start(handler)
}
