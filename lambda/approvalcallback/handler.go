package main

import (
	"context"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"html"
	"log/slog"
	"net/url"

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
	bodyParams := url.Values{}
	if request.Body != "" {
		body := request.Body
		if request.IsBase64Encoded {
			if decoded, err := base64.StdEncoding.DecodeString(body); err == nil {
				body = string(decoded)
			}
		}
		bodyParams, _ = url.ParseQuery(body)
	}
	getParam := func(key string) string {
		if v := bodyParams.Get(key); v != "" {
			return v
		}
		return request.QueryStringParameters[key]
	}

	slog.Info("approvalcallback handler invoked", "action", getParam("action"))

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	// Validate shared secret if configured
	if expectedSecret := h.cfg.AWS.SF.ApprovalSecret; expectedSecret != "" {
		providedSecret := getParam("secret")
		if subtle.ConstantTimeCompare([]byte(expectedSecret), []byte(providedSecret)) != 1 {
			slog.Warn("approvalcallback rejected: invalid secret")
			return events.APIGatewayProxyResponse{
				StatusCode: 403,
				Headers:    map[string]string{"Content-Type": "text/html"},
				Body:       "<html><body><h1>Forbidden</h1><p>Invalid or missing secret.</p></body></html>",
			}, nil
		}
	}

	token := getParam("token")
	action := getParam("action")
	refinementContext := getParam("context")
	if len(refinementContext) > 500 {
		refinementContext = refinementContext[:500]
	}

	if token == "" || action == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       "<html><body><h1>Bad Request</h1><p>Missing token or action parameter.</p></body></html>",
		}, nil
	}

	err := h.usecase.Execute(ctx, token, action, refinementContext)
	if err != nil {
		slog.Error("approvalcallback failed", "error", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    map[string]string{"Content-Type": "text/html"},
			Body:       fmt.Sprintf("<html><body><h1>Error</h1><p>%s</p></body></html>", html.EscapeString(err.Error())),
		}, nil
	}

	slog.Info("approvalcallback succeeded", "action", action)

	var message string
	switch action {
	case "approve":
		message = "Approved! The message will be sent shortly."
	case "regenerate":
		message = "Regenerating a fresh thank-you message. A new approval email will arrive shortly."
	case "refine":
		message = "Regenerating with your context. A new approval email will arrive shortly."
	case "reject":
		message = "Rejected. The message will not be sent."
	default:
		message = fmt.Sprintf("Action %q processed.", action)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       fmt.Sprintf("<html><body><h1>%s</h1></body></html>", message),
	}, nil
}
