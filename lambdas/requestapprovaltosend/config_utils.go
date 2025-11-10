package main

import (
	"fmt"
)

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
	}

	return nil
}
