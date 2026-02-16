package main

import (
	"context"
	"fmt"
	"net/url"

	ra "github.com/Bouzomgi/nycares-project-welcomer/internal/app/requestapproval"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type RequestApprovalHandler struct {
	usecase *ra.RequestApprovalUseCase
	cfg     *ra.Config
}

func NewRequestApprovalHandler(u *ra.RequestApprovalUseCase, cfg *ra.Config) *RequestApprovalHandler {
	return &RequestApprovalHandler{usecase: u, cfg: cfg}
}

func (h *RequestApprovalHandler) Handle(ctx context.Context, input models.RequestApprovalInput) (models.RequestApprovalOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	callbackEndpoint, err := url.Parse(h.cfg.AWS.SF.CallbackEndpoint)
	if err != nil {
		return models.RequestApprovalOutput{}, fmt.Errorf("callback url is invalid")
	}

	err = h.usecase.Execute(ctx, *callbackEndpoint, input.TaskToken)
	if err != nil {
		return models.RequestApprovalOutput{}, err
	}

	requestApprovalOutput := models.RequestApprovalOutput{
		Auth:                        input.Auth,
		ExistingProjectNotification: input.ExistingProjectNotification,
		MessageToSend:               input.MessageToSend,
	}

	return requestApprovalOutput, nil
}
