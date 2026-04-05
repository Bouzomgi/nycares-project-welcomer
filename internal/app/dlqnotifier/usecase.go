package dlqnotifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/email"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
)

type DLQNotifierUseCase struct {
	snsSrv snsservice.NotificationService
}

func NewDLQNotifierUseCase(snsSrv snsservice.NotificationService) *DLQNotifierUseCase {
	return &DLQNotifierUseCase{
		snsSrv: snsSrv,
	}
}

func (u *DLQNotifierUseCase) Execute(ctx context.Context, errorInfo models.DLQNotifierInput) error {
	// Extract errorMessage from the Cause JSON blob if possible.
	errorMessage := errorInfo.Cause
	var cause struct {
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal([]byte(errorInfo.Cause), &cause); err == nil && cause.ErrorMessage != "" {
		errorMessage = cause.ErrorMessage
	}

	subject, plainText, htmlBody := email.WorkflowFailed(errorInfo.FailedStep, errorMessage)

	_, err := u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, subject)
	if err != nil {
		return fmt.Errorf("failed to publish DLQ notification: %w", err)
	}

	return nil
}
