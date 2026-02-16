package main

import (
	"context"
	"crypto/subtle"
	"fmt"
	"html"

	ac "github.com/Bouzomgi/nycares-project-welcomer/internal/app/approvalcallback"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/aws/aws-lambda-go/events"
)

type ApprovalCallbackHandler struct {
	usecase *ac.ApprovalCallbackUseCase
	cfg     *ac.Config
}

func NewApprovalCallbackHandler(u *ac.ApprovalCallbackUseCase, cfg *ac.Config) *ApprovalCallbackHandler {
	return &ApprovalCallbackHandler{usecase: u, cfg: cfg}
}

func (h *ApprovalCallbackHandler) Handle(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	// Validate shared secret if configured
	if expectedSecret := h.cfg.AWS.SF.ApprovalSecret; expectedSecret != "" {
		providedSecret := request.QueryStringParameters["secret"]
		if subtle.ConstantTimeCompare([]byte(expectedSecret), []byte(providedSecret)) != 1 {
			return events.APIGatewayProxyResponse{
				StatusCode: 403,
				Headers:    map[string]string{"Content-Type": "text/html"},
				Body:       "<html><body><h1>Forbidden</h1><p>Invalid or missing secret.</p></body></html>",
			}, nil
		}
	}

	token := request.QueryStringParameters["token"]
	action := request.QueryStringParameters["action"]

	if token == "" || action == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<html><body><h1>Bad Request</h1><p>Missing token or action parameter.</p></body></html>",
		}, nil
	}

	approved := action == "approve"

	err := h.usecase.Execute(ctx, token, approved)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       fmt.Sprintf("<html><body><h1>Error</h1><p>%s</p></body></html>", html.EscapeString(err.Error())),
		}, nil
	}

	var message string
	if approved {
		message = "Approved! The message will be sent shortly."
	} else {
		message = "Rejected. The message will not be sent."
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       fmt.Sprintf("<html><body><h1>%s</h1></body></html>", message),
	}, nil
}
