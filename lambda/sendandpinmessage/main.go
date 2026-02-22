package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	spm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/sendandpinmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func buildHandler() (*SendAndPinMessageHandler, error) {
	cfg, err := config.LoadConfig[spm.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	baseUrl := endpoints.BaseUrl
	if cfg.Api.BaseUrl != "" {
		baseUrl = cfg.Api.BaseUrl
	}

	httpSvc, err := httpservice.NewHttpService(baseUrl)
	if err != nil {
		return nil, err
	}

	s3Opts := []func(*s3.Options){}
	if cfg.AWS.S3.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.AWS.S3.Endpoint)
			o.UsePathStyle = true
		})
	}
	s3Client := s3.NewFromConfig(awsCfg, s3Opts...)
	s3Svc := s3service.NewS3Service(s3Client, cfg.AWS.S3.BucketName)

	usecase := spm.NewSendAndPinMessageUseCase(s3Svc, httpSvc)
	return NewSendAndPinMessageHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		var input models.SendAndPinMessageInput
		if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
			panic(fmt.Errorf("failed to decode input: %w", err))
		}

		output, err := handler.Handle(context.Background(), input)
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
