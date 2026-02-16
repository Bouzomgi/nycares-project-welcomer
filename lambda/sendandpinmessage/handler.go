package main

import (
	"context"
	"log/slog"

	spm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/sendandpinmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type SendAndPinMessageHandler struct {
	usecase *spm.SendAndPinMessageUseCase
	cfg     *spm.Config
}

func NewSendAndPinMessageHandler(u *spm.SendAndPinMessageUseCase, cfg *spm.Config) *SendAndPinMessageHandler {
	return &SendAndPinMessageHandler{usecase: u, cfg: cfg}
}

func (h *SendAndPinMessageHandler) Handle(ctx context.Context, input models.SendAndPinMessageInput) (models.SendAndPinMessageOutput, error) {
	slog.Info("sendandpinmessage handler invoked", "projectId", input.ExistingProjectNotification.Id)

	ctx, cancel := context.WithTimeout(ctx, config.HTTPHandlerTimeout)
	defer cancel()

	auth := models.ConvertAuth(input.Auth)

	err := h.usecase.Execute(ctx, auth, input.ExistingProjectNotification.Id, input.MessageToSend.TemplateRef)
	if err != nil {
		slog.Error("sendandpinmessage failed", "error", err)
		return models.SendAndPinMessageOutput{}, err
	}

	slog.Info("sendandpinmessage succeeded")
	sendAndPinMessageOutput := models.SendAndPinMessageOutput(input)

	return sendAndPinMessageOutput, nil
}
