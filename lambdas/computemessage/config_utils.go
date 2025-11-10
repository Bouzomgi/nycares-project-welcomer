package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// --- AWS Client Setup ---

func InitAWSClients(cfg *Config, isLocal bool) (*dynamodb.Client, error) {
	var awsCfg aws.Config
	var err error

	if isLocal {
		awsCfg, err = config.LoadDefaultConfig(context.Background(),
			config.WithRegion(cfg.AWS.Dynamo.Region),
			config.WithCredentialsProvider(
				credentials.NewStaticCredentialsProvider(
					cfg.AWS.Credentials.AccessKeyID,
					cfg.AWS.Credentials.SecretAccessKey,
					"",
				),
			),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS SDK config for local: %w", err)
		}

		// DynamoDB local endpoint
		endpoint := cfg.AWS.Dynamo.Endpoint
		if endpoint == "" {
			endpoint = "http://localhost:8000"
		}

		return dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
			o.BaseEndpoint = &endpoint
		}), nil
	}

	// Lambda: use IAM role automatically
	awsCfg, err = config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Credentials.Region),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config for Lambda: %w", err)
	}

	return dynamodb.NewFromConfig(awsCfg), nil
}

// --- Config Validation ---

func ValidateConfig(cfg *Config, isLocal bool) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}

	if cfg.AWS.Credentials.Region == "" {
		return fmt.Errorf("AWS region is missing")
	}

	if isLocal {
		if cfg.AWS.Credentials.AccessKeyID == "" || cfg.AWS.Credentials.SecretAccessKey == "" {
			return fmt.Errorf("local AWS credentials missing")
		}
		if cfg.AWS.Dynamo.Endpoint == "" {
			log.Warn("DynamoDB local endpoint not set, using default http://localhost:8000")
		}
		if cfg.AWS.S3.Endpoint == "" {
			log.Warn("S3 local endpoint not set")
		}
	}

	return nil
}
