package dlqnotifier

import (
	"context"
	"encoding/json"
	"fmt"
	"html"

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
	subject := "NYC Cares Project Welcomer — Workflow Step Failed"

	// Extract errorMessage from the Cause JSON blob if possible.
	errorMessage := errorInfo.Cause
	var cause struct {
		ErrorMessage string `json:"errorMessage"`
	}
	if err := json.Unmarshal([]byte(errorInfo.Cause), &cause); err == nil && cause.ErrorMessage != "" {
		errorMessage = cause.ErrorMessage
	}

	plainText := fmt.Sprintf("Workflow step failed.\nStep: %s\nError: %s", errorInfo.FailedStep, errorMessage)

	htmlBody := fmt.Sprintf(`<h2>Workflow Step Failed</h2>
<table>
  <tr><td><strong>Step</strong></td><td>%s</td></tr>
  <tr><td><strong>Error</strong></td><td>%s</td></tr>
</table>`,
		html.EscapeString(errorInfo.FailedStep),
		html.EscapeString(errorMessage),
	)

	_, err := u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, subject)
	if err != nil {
		return fmt.Errorf("failed to publish DLQ notification: %w", err)
	}

	return nil
}
