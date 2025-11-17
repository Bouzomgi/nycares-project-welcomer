package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// LoadConfig is a unified config loader that works with any struct type.
// It supports environment-specific configs and Lambda deployments
//
// Parameters:
// 		- env: The environment to load config for ("dev", "prod", etc.)
//					 If empty, defaults to "local"

func LoadConfig(env string) (*Config, error) {
	v := viper.New()

	// Set default environment if not specified
	if env == "" {
		env = "local"
	}

	// Set up paths
	configDir := "configs"
	configName := "config"

	// Check if running in Lambda
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		configDir = "/opt/config" // Lambda layer config path
	}

	configPath := filepath.Join(configDir, env)

	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// Environment variables take precedence
	v.AutomaticEnv()
	v.SetEnvPrefix("NYCARES")

	// Load the config file
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set the environment in the config\
	config.Env = env

	return &config, nil
}
