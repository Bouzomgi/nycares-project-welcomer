package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"nycaresprojectwelcomer/internal/models"

	"github.com/aws/aws-lambda-go/lambda"
)

// LambdaOutput defines what the Lambda returns.
type LambdaOutput struct {
	Cookies   map[string]string `json:"cookies"`
	Schedule []models.Project `json:"schedule"`
}

// handler runs on Lambda invocation.
func handler(ctx context.Context) (LambdaOutput, error) {

	cfg, err := LoadConfig()
	if err != nil {
			log.Fatal("failed to load config:", err)
	}

	client := &http.Client{}

	creds := Credentials {
		Username: cfg.Account.Username,
		Password: cfg.Account.Password,
	}

	log.Print("Attempting login")
	cookies, err := Login(client, PostLoginUrl, creds)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return LambdaOutput{}, err
	}

	log.Print("Fetching schedule")
	schedule, err := GetSchedule(client, cfg.Account.InternalId, cookies)
	if err != nil {
		log.Printf("GetSchedule failed: %v", err)
		return LambdaOutput{}, err
	}

	log.Printf("Fetched schedule, has %d events", len(schedule))

	return LambdaOutput{
		Cookies:   cookies,
		Schedule: schedule,
	}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		// Local debug
		resp, _ := handler(context.Background())
		fmt.Println(resp)
		return
	}

	// Lambda runtime
	lambda.Start(handler)
}