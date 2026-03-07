package config

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// loadSecretsToEnv fetches secrets from AWS Secrets Manager and injects them
// as environment variables so the existing viper config loading picks them up.
// Only runs when NYCARES_SECRET_ARN is set (i.e. in AWS Lambda, not locally).
func loadSecretsToEnv() error {
	secretArn := os.Getenv("NYCARES_SECRET_ARN")
	if secretArn == "" {
		return nil
	}

	ctx := context.Background()
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config for secrets: %w", err)
	}

	client := secretsmanager.NewFromConfig(cfg)
	result, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretArn,
	})
	if err != nil {
		return fmt.Errorf("failed to get secret %s: %w", secretArn, err)
	}

	if result.SecretString == nil {
		return fmt.Errorf("secret %s has no string value", secretArn)
	}

	var secrets map[string]string
	if err := json.Unmarshal([]byte(*result.SecretString), &secrets); err != nil {
		return fmt.Errorf("failed to parse secret JSON: %w", err)
	}

	for key, val := range secrets {
		os.Setenv(key, val)
	}

	return nil
}
