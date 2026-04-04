package notifycompletion

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/email"
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

	projectDate := utils.DateToString(project.Date)
	subject, plainText, htmlBody := email.Completion(notificationType.String(), project.Name, projectDate)

	_, err := u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, subject)
	if err != nil {
		return fmt.Errorf("failed to publish completion notification: %w", err)
	}

	return nil
}
