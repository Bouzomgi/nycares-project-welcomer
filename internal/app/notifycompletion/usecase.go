package notifycompletion

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type NotifyCompletionUseCase struct {
	snsSrv snsservice.NotificationService
}

func NewNotifyCompletionUseCase(snsSrv snsservice.NotificationService) *NotifyCompletionUseCase {
	return &NotifyCompletionUseCase{
		snsSrv: snsSrv,
	}
}

func (u *NotifyCompletionUseCase) Execute(
	ctx context.Context,
	notificationType domain.NotificationType,
	project domain.Project,
) error {

	completionMessage := createCompletionMessage(notificationType, project)
	subject := "Message Sent!"

	_, err := u.snsSrv.PublishNotification(ctx, completionMessage, subject)
	if err != nil {
		return fmt.Errorf("failed to publish completion notification: %w", err)
	}

	return nil
}

func createCompletionMessage(messageType domain.NotificationType, project domain.Project) string {
	projectDate := utils.DateToString(project.Date)
	return fmt.Sprintf(
		"Successfully sent %s message to %s on %s!",
		messageType.String(),
		project.Name,
		projectDate,
	)
}
