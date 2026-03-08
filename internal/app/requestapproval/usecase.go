package requestapproval

import (
	"context"
	"fmt"
	"net/url"

	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
)

type RequestApprovalUseCase struct {
	snsSrv         snsservice.NotificationService
	approvalSecret string
}

func NewRequestApprovalUseCase(snsSrv snsservice.NotificationService, approvalSecret string) *RequestApprovalUseCase {
	return &RequestApprovalUseCase{
		snsSrv:         snsSrv,
		approvalSecret: approvalSecret,
	}
}

func (u *RequestApprovalUseCase) Execute(ctx context.Context, callbackEndpoint url.URL, taskToken string) error {

	if taskToken == "" {
		return fmt.Errorf("taskToken must be defined")
	}

	approveLink := buildCallbackLink(callbackEndpoint, taskToken, "approve", u.approvalSecret)
	rejectLink := buildCallbackLink(callbackEndpoint, taskToken, "reject", u.approvalSecret)

	plainText := fmt.Sprintf("Approve: %s\n\nReject: %s", approveLink, rejectLink)
	htmlBody := fmt.Sprintf(
		`<p><a href="%s">Approve</a></p><p><a href="%s">Reject</a></p>`,
		approveLink, rejectLink,
	)

	_, err := u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, "Project Message Approval")
	if err != nil {
		return fmt.Errorf("failed to publish approval notification: %w", err)
	}

	return nil
}

func buildCallbackLink(baseURL url.URL, taskToken string, action string, secret string) string {

	// Append "/callback" to path correctly
	if len(baseURL.Path) == 0 || baseURL.Path[len(baseURL.Path)-1] != '/' {
		baseURL.Path += "/callback"
	} else {
		baseURL.Path += "callback"
	}

	// Add the task token, action, and secret as query parameters
	q := baseURL.Query()
	q.Set("token", taskToken)
	q.Set("action", action)
	if secret != "" {
		q.Set("secret", secret)
	}
	baseURL.RawQuery = q.Encode()

	return baseURL.String()
}
