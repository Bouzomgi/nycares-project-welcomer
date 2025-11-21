package notifycompletion

import (
	"context"
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type NotifyCompletionUseCase struct {
	snsSrv *snsservice.SNSService
}

func NewNotifyCompletionUseCase(snsSrv *snsservice.SNSService) *NotifyCompletionUseCase {
	return &NotifyCompletionUseCase{
		snsSrv: snsSrv,
	}
}

func (u *NotifyCompletionUseCase) Execute(ctx context.Context, notificationType domain.NotificationType, project domain.Project) error {

	completionMessage := createCompletionMessage(notificationType, project)

	subject := "Message Sent!"
	u.snsSrv.PublishMessage(ctx, completionMessage, subject)

	return nil
}

func createCompletionMessage(messageType domain.NotificationType, project domain.Project) string {
	projectDateString := utils.DateToString(project.Date)

	return fmt.Sprintf("Successfully sent %s message to %s on %s!", messageType, project.Name, projectDateString)
}
