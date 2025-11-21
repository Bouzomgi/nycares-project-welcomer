package main

import (
	"context"
	"time"

	cm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/computemessage"
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

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	domainProject, err := models.BuildDomainProject(input.Project)
	if err != nil {
		return models.ComputeMessageOutput{}, err
	}

	messageType, messageRef, err := h.usecase.Execute(ctx, h.cfg.AWS.S3.BucketName, domainProject)
	if err != nil {
		return models.ComputeMessageOutput{}, err
	}

	message := models.BuildMessage(messageType.String(), messageRef)

	output := models.ComputeMessageOutput{
		Auth:          input.Auth,
		Project:       input.Project,
		MessageToSend: message,
	}

	return output, nil
}
