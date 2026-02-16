package main

import (
	"context"
	"log/slog"

	cm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/computemessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
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
	slog.Info("computemessage handler invoked", "project", input.Project)

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	domainProject, err := models.BuildDomainProject(input.Project)
	if err != nil {
		slog.Error("computemessage failed to build project", "error", err)
		return models.ComputeMessageOutput{}, err
	}

	existingNotification, messageType, messageRef, err := h.usecase.Execute(ctx, h.cfg.AWS.S3.BucketName, domainProject)
	if err != nil {
		slog.Error("computemessage failed", "error", err)
		return models.ComputeMessageOutput{}, err
	}

	slog.Info("computemessage succeeded", "messageType", messageType.String())

	message := models.BuildMessage(messageType.String(), messageRef)
	outputNotification := models.ConvertDomainProjectNotification(existingNotification)

	output := models.ComputeMessageOutput{
		Auth:                        input.Auth,
		ExistingProjectNotification: outputNotification,
		MessageToSend:               message,
	}

	return output, nil
}
