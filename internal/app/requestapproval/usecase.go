package requestapproval

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/email"
	s3service "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/s3"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
)

type RequestApprovalUseCase struct {
	snsSrv         snsservice.NotificationService
	s3Srv          s3service.ContentService
	approvalSecret string
}

func NewRequestApprovalUseCase(snsSrv snsservice.NotificationService, s3Srv s3service.ContentService, approvalSecret string) *RequestApprovalUseCase {
	return &RequestApprovalUseCase{
		snsSrv:         snsSrv,
		s3Srv:          s3Srv,
		approvalSecret: approvalSecret,
	}
}

func (u *RequestApprovalUseCase) Execute(ctx context.Context, callbackEndpoint url.URL, taskToken, projectName, projectDate, messageType, templateRef string, mockMode bool) error {

	if taskToken == "" {
		return fmt.Errorf("taskToken must be defined")
	}

	messageContent, err := u.s3Srv.GetMessageContent(ctx, templateRef)
	if err != nil {
		return fmt.Errorf("failed to fetch message content: %w", err)
	}

	messageContent = strings.ReplaceAll(messageContent, "{{projectName}}", projectName)

	approveLink := buildCallbackLink(callbackEndpoint, taskToken, "approve", u.approvalSecret)
	rejectLink := buildCallbackLink(callbackEndpoint, taskToken, "reject", u.approvalSecret)

	subject, plainText, htmlBody := email.ApprovalRequest(projectName, projectDate, messageType, messageContent, approveLink, rejectLink, mockMode)

	_, err = u.snsSrv.PublishHTMLEmailNotification(ctx, plainText, htmlBody, subject)
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
