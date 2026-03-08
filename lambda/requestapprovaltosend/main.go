package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	ra "github.com/Bouzomgi/nycares-project-welcomer/internal/app/requestapproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func buildHandler() (*RequestApprovalHandler, error) {
	cfg, err := config.LoadConfig[ra.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	snsClient := sns.NewFromConfig(awsCfg)
	snsSvc := snsservice.NewSNSService(snsClient, cfg.AWS.SNS.TopicArn)

	s3Opts := []func(*s3.Options){}
	if cfg.AWS.S3.Endpoint != "" {
		s3Opts = append(s3Opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.AWS.S3.Endpoint)
			o.UsePathStyle = true
		})
	}
	s3Client := s3.NewFromConfig(awsCfg, s3Opts...)
	s3Svc := s3service.NewS3Service(s3Client, cfg.AWS.S3.BucketName)

	usecase := ra.NewRequestApprovalUseCase(snsSvc, s3Svc, cfg.AWS.SF.ApprovalSecret)
	return NewRequestApprovalHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		output, err := handler.Handle(context.Background(), models.RequestApprovalInput{})
		if err != nil {
			panic(err)
		}
		data, _ := json.MarshalIndent(output, "", "  ")
		fmt.Println(string(data))
		return
	}

	lambda.Start(handler.Handle)
}
