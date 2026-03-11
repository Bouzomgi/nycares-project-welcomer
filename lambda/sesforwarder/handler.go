package main

import (
	"context"
	"encoding/json"
	"log/slog"

	sesservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/ses"
	"github.com/aws/aws-lambda-go/events"
)

type SESForwarderHandler struct {
	sesSvc *sesservice.SESService
}

func NewSESForwarderHandler(sesSvc *sesservice.SESService) *SESForwarderHandler {
	return &SESForwarderHandler{sesSvc: sesSvc}
}

type htmlPayload struct {
	HTMLBody  string `json:"htmlBody"`
	PlainText string `json:"plainText"`
}

func (h *SESForwarderHandler) Handle(ctx context.Context, event events.SNSEvent) error {
	for _, record := range event.Records {
		msg := record.SNS
		subject := msg.Subject
		if subject == "" {
			subject = "(no subject)"
		}

		attr, hasFormat := msg.MessageAttributes["format"]
		if hasFormat && attr.(map[string]interface{})["Value"] == "html" {
			var payload htmlPayload
			if err := json.Unmarshal([]byte(msg.Message), &payload); err != nil {
				slog.Error("sesforwarder failed to parse HTML payload", "error", err)
				return err
			}
			if err := h.sesSvc.SendHTMLEmail(ctx, subject, payload.HTMLBody, payload.PlainText); err != nil {
				slog.Error("sesforwarder failed to send HTML email", "error", err)
				return err
			}
		} else {
			if err := h.sesSvc.SendPlainEmail(ctx, subject, msg.Message); err != nil {
				slog.Error("sesforwarder failed to send plain email", "error", err)
				return err
			}
		}

		slog.Info("sesforwarder sent email", "subject", subject, "messageId", msg.MessageID)
	}
	return nil
}
