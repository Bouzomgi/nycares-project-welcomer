package config

import "time"

const (
	// DefaultHandlerTimeout is the timeout for handlers that perform
	// internal-only operations (DynamoDB, SNS, S3).
	DefaultHandlerTimeout = 10 * time.Second

	// HTTPHandlerTimeout is the timeout for handlers that make external
	// HTTP calls (login, fetch projects, send messages).
	// Set below the Lambda function timeout (30s) to allow graceful error propagation.
	HTTPHandlerTimeout = 25 * time.Second
)
