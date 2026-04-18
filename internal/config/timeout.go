package config

import (
	"context"
	"time"
)

const (
	// DefaultHandlerTimeout is the timeout for handlers that perform
	// internal-only operations (DynamoDB, SNS, S3).
	DefaultHandlerTimeout = 10 * time.Second

	// HTTPHandlerTimeout is the timeout for handlers that make external
	// HTTP calls (login, fetch projects, send messages).
	// Set below the Lambda function timeout (30s) to allow graceful error propagation.
	HTTPHandlerTimeout = 25 * time.Second

	// AIHandlerTimeout is the timeout for handlers that call AI generation
	// services (e.g. Bedrock). Set below the Lambda function timeout (60s).
	AIHandlerTimeout = 55 * time.Second
)

// HandlerDeadline returns a context capped at the earlier of the given timeout
// and the Lambda context's remaining execution time, preventing the handler
// from outlasting the Lambda's hard kill deadline.
func HandlerDeadline(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	deadline := time.Now().Add(timeout)
	if lambdaDeadline, ok := ctx.Deadline(); ok && lambdaDeadline.Before(deadline) {
		deadline = lambdaDeadline
	}
	return context.WithDeadline(ctx, deadline)
}
