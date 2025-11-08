package main

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/confighelper"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type LambdaOutput struct {
	Auth models.Auth `json:"auth"`
}

func createAuthFromCookies(cookies []*http.Cookie) models.Auth {
	var auth models.Auth
	for _, c := range cookies {
		auth.Cookies = append(auth.Cookies, models.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Domain: c.Domain,
			Path:   c.Path,
		})
	}
	return auth
}

func handler(ctx context.Context) (LambdaOutput, error) {

	cfg, err := confighelper.LoadConfig[Config]()
	if err != nil {
		log.Errorf("failed to load config: %v", err)
		return LambdaOutput{}, err
	}

	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	defer client.CloseIdleConnections()

	creds := Credentials{
		Username: cfg.Account.Username,
		Password: cfg.Account.Password,
	}

	log.Info("Attempting login")
	cookies, err := Login(client, endpoints.BaseUrl, creds)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		return LambdaOutput{}, err
	}
	log.Info("Login succeeded")

	auth := createAuthFromCookies(cookies)

	return LambdaOutput{
		Auth: auth,
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
