package requestapproval

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type RequestApprovalHandler struct {
	usecase *RequestApprovalUseCase
	cfg     *Config
}

func NewRequestApprovalHandler(u *RequestApprovalUseCase, cfg *Config) *RequestApprovalHandler {
	return &RequestApprovalHandler{usecase: u, cfg: cfg}
}

func (h *RequestApprovalHandler) Handle(ctx context.Context, input models.RequestApprovalInput) (models.RequestApprovalOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	callbackEndpoint, err := url.Parse(h.cfg.CallbackEndpoint)
	if err != nil {
		return models.RequestApprovalOutput{}, fmt.Errorf("callback url is invalid")
	}

	err = h.usecase.Execute(ctx, *callbackEndpoint, input.TaskToken)
	if err != nil {
		return models.RequestApprovalOutput{}, err
	}

	requestApprovalOutput := models.RequestApprovalOutput{
		Auth:          input.Auth,
		Project:       input.Project,
		MessageToSend: input.MessageToSend,
	}

	return requestApprovalOutput, nil
}
