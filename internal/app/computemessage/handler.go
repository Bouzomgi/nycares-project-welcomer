package computemessage

import (
	"context"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type ComputeMessageHandler struct {
	usecase *ComputeMessageUseCase
	cfg     *Config
}

func NewComputeMessageHandler(u *ComputeMessageUseCase, cfg *Config) *ComputeMessageHandler {
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
