package config

import "time"

const (
	// DefaultHandlerTimeout is the timeout for handlers that perform
	// internal-only operations (DynamoDB, SNS, S3).
	DefaultHandlerTimeout = 10 * time.Second

	// HTTPHandlerTimeout is the timeout for handlers that make external
	// HTTP calls (login, fetch projects, send messages).
	HTTPHandlerTimeout = 30 * time.Second
)
