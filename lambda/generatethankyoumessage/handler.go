package main

import (
	"context"
	"log/slog"

	gtm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/generatethankyoumessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type GenerateThankYouMessageHandler struct {
	usecase *gtm.GenerateThankYouMessageUseCase
	cfg     *gtm.Config
}

func NewGenerateThankYouMessageHandler(u *gtm.GenerateThankYouMessageUseCase, cfg *gtm.Config) *GenerateThankYouMessageHandler {
	return &GenerateThankYouMessageHandler{usecase: u, cfg: cfg}
}

func (h *GenerateThankYouMessageHandler) Handle(ctx context.Context, input models.GenerateThankYouMessageInput) (models.GenerateThankYouMessageOutput, error) {
	slog.Info("generatethankyoumessage handler invoked", "executionId", input.ExecutionId, "project", input.ExistingProjectNotification.Name)

	ctx, cancel := context.WithTimeout(ctx, config.AIHandlerTimeout)
	defer cancel()

	generatedContent, err := h.usecase.Execute(ctx, input.ExistingProjectNotification.Name)
	if err != nil {
		slog.Error("generatethankyoumessage failed", "executionId", input.ExecutionId, "error", err)
		return models.GenerateThankYouMessageOutput{}, err
	}

	msg := models.BuildMessage("thankYou", "")
	msg.GeneratedContent = generatedContent

	slog.Info("generatethankyoumessage succeeded", "executionId", input.ExecutionId)

	return models.GenerateThankYouMessageOutput{
		Auth:                        input.Auth,
		ExistingProjectNotification: input.ExistingProjectNotification,
		MessageToSend:               msg,
		TargetSendTime:              input.TargetSendTime,
		ExecutionId:                 input.ExecutionId,
	}, nil
}
