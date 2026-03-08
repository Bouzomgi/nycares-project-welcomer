package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// loadSecretsToEnv fetches parameters from AWS SSM Parameter Store and injects
// them as environment variables so the existing viper config loading picks them up.
// Only runs when NYCARES_SSM_PATH is set (i.e. in AWS Lambda, not locally).
func loadSecretsToEnv() error {
	ssmPath := os.Getenv("NYCARES_SSM_PATH")
	if ssmPath == "" {
		return nil
	}

	ctx := context.Background()
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config for SSM: %w", err)
	}

	client := ssm.NewFromConfig(cfg)
	result, err := client.GetParametersByPath(ctx, &ssm.GetParametersByPathInput{
		Path:           aws.String(ssmPath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to get SSM parameters at %s: %w", ssmPath, err)
	}

	for _, p := range result.Parameters {
		key := strings.TrimPrefix(*p.Name, ssmPath)
		os.Setenv(key, *p.Value)
	}

	return nil
}
