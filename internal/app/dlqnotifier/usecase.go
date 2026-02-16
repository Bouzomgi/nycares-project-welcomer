package dlqnotifier

import (
	"context"
	"fmt"

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
	message := fmt.Sprintf("Workflow step failed.\nError: %s\nCause: %s", errorInfo.Error, errorInfo.Cause)
	subject := "Workflow Error"

	_, err := u.snsSrv.PublishNotification(ctx, message, subject)
	if err != nil {
		return fmt.Errorf("failed to publish DLQ notification: %w", err)
	}

	return nil
}
