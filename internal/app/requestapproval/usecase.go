package requestapproval

import (
	"context"
	"fmt"
	"net/url"

	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
)

type RequestApprovalUseCase struct {
	snsSrv snsservice.NotificationService
}

func NewRequestApprovalUseCase(snsSrv snsservice.NotificationService) *RequestApprovalUseCase {
	return &RequestApprovalUseCase{
		snsSrv: snsSrv,
	}
}

func (u *RequestApprovalUseCase) Execute(ctx context.Context, callbackEndpoint url.URL, taskToken string) error {

	if taskToken == "" {
		return fmt.Errorf("taskToken must be defined")
	}

	approveLink := buildCallbackLink(callbackEndpoint, taskToken, "approve")
	rejectLink := buildCallbackLink(callbackEndpoint, taskToken, "reject")

	message := fmt.Sprintf("Approve: %s\n\nReject: %s", approveLink, rejectLink)

	u.snsSrv.PublishNotification(ctx, message, "Project Message Approval")

	return nil
}

func buildCallbackLink(baseURL url.URL, taskToken string, action string) string {

	// Append "/callback" to path correctly
	if len(baseURL.Path) == 0 || baseURL.Path[len(baseURL.Path)-1] != '/' {
		baseURL.Path += "/callback"
	} else {
		baseURL.Path += "callback"
	}

	// Add the task token and action as query parameters
	q := baseURL.Query()
	q.Set("token", taskToken)
	q.Set("action", action)
	baseURL.RawQuery = q.Encode()

	return baseURL.String()
}
