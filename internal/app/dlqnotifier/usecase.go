package dlqnotifier

import (
	"context"
	"fmt"

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

func (u *DLQNotifierUseCase) Execute(ctx context.Context, errorInfo map[string]interface{}) error {
	message := fmt.Sprintf("Workflow step failed with error: %v", errorInfo)
	subject := "Workflow Error"

	_, err := u.snsSrv.PublishNotification(ctx, message, subject)
	if err != nil {
		return fmt.Errorf("failed to publish DLQ notification: %w", err)
	}

	return nil
}
