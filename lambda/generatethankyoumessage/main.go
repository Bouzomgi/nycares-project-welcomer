package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	gtm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/generatethankyoumessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/platform/awsconfig"
	bedrockservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/bedrock"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func buildHandler() (*GenerateThankYouMessageHandler, error) {
	cfg, err := config.LoadConfig[gtm.Config]()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	awsCfg, err := awsconfig.LoadAWSConfigFromConfig(ctx, cfg)
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

	var bedrockSvc bedrockservice.GenerationService
	if cfg.Mock.GenerateThankYou {
		bedrockSvc = bedrockservice.NewMockBedrockService()
	} else {
		bedrockOpts := []func(*bedrockruntime.Options){}
		if cfg.AWS.Bedrock.Endpoint != "" {
			bedrockOpts = append(bedrockOpts, func(o *bedrockruntime.Options) {
				o.BaseEndpoint = aws.String(cfg.AWS.Bedrock.Endpoint)
			})
		}
		bedrockClient := bedrockruntime.NewFromConfig(awsCfg, bedrockOpts...)
		bedrockSvc = bedrockservice.NewBedrockService(bedrockClient)
	}

	usecase := gtm.NewGenerateThankYouMessageUseCase(s3Svc, bedrockSvc, cfg.AWS.S3.BucketName)
	return NewGenerateThankYouMessageHandler(usecase, cfg), nil
}

func main() {
	config.InitLogging()
	handler, err := buildHandler()
	if err != nil {
		panic(err)
	}

	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		var input models.GenerateThankYouMessageInput
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
