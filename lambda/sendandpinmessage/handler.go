package main

import (
	"context"
	"time"

	spm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/sendandpinmessage"
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

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	auth := models.ConvertAuth(input.Auth)

	err := h.usecase.Execute(ctx, auth, input.Project.Id, input.MessageToSend.TemplateRef)
	if err != nil {
		return models.SendAndPinMessageOutput{}, err
	}

	sendAndPinMessageOutput := models.SendAndPinMessageOutput(input)

	return sendAndPinMessageOutput, nil
}
