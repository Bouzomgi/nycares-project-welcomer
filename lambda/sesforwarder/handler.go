package main

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/email"
	sesservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/ses"
	snsservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/sns"
	"github.com/aws/aws-lambda-go/events"
)

type SESForwarderHandler struct {
	sesSvc *sesservice.SESService
	snsSvc snsservice.NotificationService
}

func NewSESForwarderHandler(sesSvc *sesservice.SESService, snsSvc snsservice.NotificationService) *SESForwarderHandler {
	return &SESForwarderHandler{sesSvc: sesSvc, snsSvc: snsSvc}
}

type htmlPayload struct {
	HTMLBody  string `json:"htmlBody"`
	PlainText string `json:"plainText"`
}

func (h *SESForwarderHandler) publishError(err error) {
	if h.snsSvc == nil {
		return
	}
	subject, plainText, htmlBody, renderErr := email.WorkflowFailed("SESForwarder", err.Error())
	if renderErr != nil {
		slog.Error("sesforwarder failed to render error notification", "error", renderErr)
		return
	}
	if _, publishErr := h.snsSvc.PublishHTMLEmailNotification(context.Background(), plainText, htmlBody, subject); publishErr != nil {
		slog.Error("sesforwarder failed to publish error notification", "error", publishErr)
	}
}

func (h *SESForwarderHandler) Handle(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		msg := record.SNS
		subject := msg.Subject
		if subject == "" {
			subject = "(no subject)"
		}

		attr, hasFormat := msg.MessageAttributes["format"]
		attrMap, attrIsMap := attr.(map[string]interface{})
		if hasFormat && attrIsMap && attrMap["Value"] == "html" {
			var payload htmlPayload
			if err := json.Unmarshal([]byte(msg.Message), &payload); err != nil {
				slog.Error("sesforwarder failed to parse HTML payload", "error", err)
				h.publishError(err)
				return err
			}
			if err := h.sesSvc.SendHTMLEmail(ctx, subject, payload.HTMLBody, payload.PlainText); err != nil {
				slog.Error("sesforwarder failed to send HTML email", "error", err)
				h.publishError(err)
				return err
			}
		} else {
			if err := h.sesSvc.SendPlainEmail(ctx, subject, msg.Message); err != nil {
				slog.Error("sesforwarder failed to send plain email", "error", err)
				h.publishError(err)
				return err
			}
		}

		slog.Info("sesforwarder sent email", "subject", subject, "messageId", msg.MessageID)
	}
	return nil
}
