package fetchprojects

import (
	"context"
	"fmt"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/httpservice"
	"github.com/aws/aws-lambda-go/lambda"

	log "github.com/sirupsen/logrus"
)

type LambdaInput struct {
	Auth models.Auth `json:"auth"`
}

type LambdaOutput struct {
	Auth     models.Auth      `json:"auth"`
	Projects []models.Project `json:"projects"`
}

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfig(os.Getenv("NYCARES_ENV"))
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
}

func handler(ctx context.Context, input LambdaInput) (LambdaOutput, error) {

	httpService, err := httpservice.NewHttpService(endpoints.BaseUrl, httpservice.WithAuth(input.Auth))
	if err != nil {
		return LambdaOutput{}, fmt.Errorf("failed to create HTTP client: %w", err)
	}

	log.Info("Fetching schedule")
	schedule, err := httpService.GetSchedule(ctx, cfg.Account.InternalId)
	if err != nil {
		return LambdaOutput{}, fmt.Errorf("failed to get schedule: %w", err)
	}

	log.Infof("Fetched schedule, has %d events", len(schedule))

	projects := models.ReduceToProjects(schedule)

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

	lambda.Start(handler)
}
