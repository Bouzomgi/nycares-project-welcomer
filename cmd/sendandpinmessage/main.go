package sendandpinmessage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/sendandpinmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func buildHandler() (*sendandpinmessage.SendAndPinMessageHandler, error) {
	cfg, err := config.LoadConfig[sendandpinmessage.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	httpSvc, err := httpservice.NewHttpService(endpoints.BaseUrl)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCfg)
	s3Svc := s3service.NewS3Service(s3Client, cfg.AWS.S3.BucketName)

	usecase := sendandpinmessage.NewSendAndPinMessageUseCase(s3Svc, httpSvc)
	return sendandpinmessage.NewSendAndPinMessageHandler(usecase, cfg), nil
}

func main() {
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("_LAMBDA_SERVER_PORT") == "" {
		output, err := handler.Handle(context.Background(), models.SendAndPinMessageInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
