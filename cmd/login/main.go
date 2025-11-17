package login

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/httpservice"

	"github.com/aws/aws-lambda-go/lambda"
	log "github.com/sirupsen/logrus"
)

type LambdaOutput struct {
	Auth models.Auth `json:"auth"`
}

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig(os.Getenv("NYCARES_ENV"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func handler(ctx context.Context) (LambdaOutput, error) {
	httpService, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return LambdaOutput{}, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	creds := models.Credentials{
		Username: cfg.Account.Username,
		Password: cfg.Account.Password,
	}

	log.WithField("env", cfg.Env).Info("Attempting login")

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	authResp, err := httpService.Login(ctx, creds)
	if err != nil {
		log.Errorf("Login failed: %v", err)
		return LambdaOutput{}, fmt.Errorf("login handler: %w", err)
	}
	log.Info("Login succeeded")

	return LambdaOutput{Auth: authResp}, nil
}

func main() {
	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		if _, err := handler(context.Background()); err != nil {
			log.Fatal(err)
		}
		return
	}

	lambda.Start(handler)
}
