package main

import (
	"context"

	rm "github.com/Bouzomgi/nycares-project-welcomer/internal/app/recordmessage"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type RecordMessageHandler struct {
	usecase *rm.RecordMessageUseCase
	cfg     *rm.Config
}

func NewRecordMessageHandler(u *rm.RecordMessageUseCase, cfg *rm.Config) *RecordMessageHandler {
	return &RecordMessageHandler{usecase: u, cfg: cfg}
}

func (h *RecordMessageHandler) Handle(ctx context.Context, input models.RecordMessageInput) (models.RecordMessageOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	domainProjectNotification, err := models.ConvertModelProjectNotification(input.ExistingProjectNotification)
	if err != nil {
		return models.RecordMessageOutput{}, err
	}

	notificationType, err := domain.ParseNotificationType(input.MessageToSend.Type)
	if err != nil {
		return models.RecordMessageOutput{}, err
	}

	updatedProjectNotification, err := h.usecase.Execute(ctx, domainProjectNotification, notificationType)
	if err != nil {
		return models.RecordMessageOutput{}, err
	}

	outputProjectNotification := models.ConvertDomainProjectNotification(updatedProjectNotification)

	output := models.RecordMessageOutput{
		TaskToken:                   input.TaskToken,
		Auth:                        input.Auth,
		ExistingProjectNotification: input.ExistingProjectNotification,
		MessageToSend:               input.MessageToSend,
		RecordedProjectNotification: outputProjectNotification,
	}

	return output, nil
}
