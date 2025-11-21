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

	link := buildCallbackLink(callbackEndpoint, taskToken)

	u.snsSrv.PublishNotification(ctx, link, "yooooooooo")

	return nil
}

func buildCallbackLink(baseURL url.URL, taskToken string) string {

	// Append "/callback" to path correctly
	if len(baseURL.Path) == 0 || baseURL.Path[len(baseURL.Path)-1] != '/' {
		baseURL.Path += "/callback"
	} else {
		baseURL.Path += "callback"
	}

	// Add the task token as a query parameter
	q := baseURL.Query()
	q.Set("token", taskToken)
	baseURL.RawQuery = q.Encode()

	return baseURL.String()
}
