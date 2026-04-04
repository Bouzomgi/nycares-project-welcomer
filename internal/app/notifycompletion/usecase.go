package notifycompletion

import (
	"context"
	"fmt"
	"html"

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
	mockMode bool,
) error {

	projectDate := utils.DateToString(project.Date)
	msgType := notificationType.String()
	projectName := project.Name

	destination := "real NYC Cares platform"
	if mockMode {
		destination = "mock server"
	}

	plainText := fmt.Sprintf(
		"Successfully sent %s message to %s on %s!\n\nSending to: %s",
		msgType, projectName, projectDate, destination,
	)

	htmlBody := fmt.Sprintf(
		`<p>Successfully sent <strong>%s</strong> message to <strong>%s</strong> on %s!</p>`+
			`<p><em>Sending to: %s</em></p>`,
		html.EscapeString(msgType),
		html.EscapeString(projectName),
		html.EscapeString(projectDate),
		html.EscapeString(destination),
	)

	_, err := u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, "Message Sent!")
	if err != nil {
		return fmt.Errorf("failed to publish completion notification: %w", err)
	}

	return nil
}
