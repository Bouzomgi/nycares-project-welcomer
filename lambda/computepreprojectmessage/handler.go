package main

import (
	"context"
	"log/slog"

	cm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/computepreprojectmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type ComputeMessageHandler struct {
	usecase *cm.ComputeMessageUseCase
	cfg     *cm.Config
}

func NewComputeMessageHandler(u *cm.ComputeMessageUseCase, cfg *cm.Config) *ComputeMessageHandler {
	return &ComputeMessageHandler{usecase: u, cfg: cfg}
}

func (h *ComputeMessageHandler) Handle(ctx context.Context, input models.ComputeMessageInput) (models.ComputeMessageOutput, error) {
	slog.Info("computepreprojectmessage handler invoked", "executionId", input.ExecutionId, "messageType", input.MessageType)

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	messageType, err := domain.ParseNotificationType(input.MessageType)
	if err != nil {
		slog.Error("computepreprojectmessage invalid message type", "executionId", input.ExecutionId, "error", err)
		return models.ComputeMessageOutput{}, err
	}

	messageRef, err := h.usecase.Execute(h.cfg.AWS.S3.BucketName, input.ExistingProjectNotification.Name, messageType)
	if err != nil {
		slog.Error("computepreprojectmessage failed", "executionId", input.ExecutionId, "error", err)
		return models.ComputeMessageOutput{}, err
	}

	slog.Info("computepreprojectmessage succeeded", "executionId", input.ExecutionId, "messageType", input.MessageType)

	return models.ComputeMessageOutput{
		Auth:                        input.Auth,
		ExistingProjectNotification: input.ExistingProjectNotification,
		MessageToSend:               models.BuildMessage(input.MessageType, messageRef),
		TargetSendTime:              input.TargetSendTime,
		ExecutionId:                 input.ExecutionId,
	}, nil
}
