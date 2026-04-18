package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/app/login"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type LoginHandler struct {
	usecase *login.LoginUseCase
	cfg     *login.Config
}

func NewLoginHandler(u *login.LoginUseCase, cfg *login.Config) *LoginHandler {
	return &LoginHandler{usecase: u, cfg: cfg}
}

func (h *LoginHandler) Handle(ctx context.Context, input models.LoginInput) (models.LoginOutput, error) {
	slog.Info("login handler invoked", "executionId", input.ExecutionId)

	var execInput struct {
		MockProjects json.RawMessage `json:"mockProjects,omitempty"`
	}
	if len(input.Context) > 0 {
		if err := json.Unmarshal(input.Context, &execInput); err != nil {
			return models.LoginOutput{}, fmt.Errorf("failed to unmarshal context: %w", err)
		}
	}

	creds := domain.Credentials{
		Username:         h.cfg.Account.Username,
		Password:         h.cfg.Account.Password,
		MockProjectsJSON: execInput.MockProjects,
	}

	ctx, cancel := config.HandlerDeadline(ctx, config.HTTPHandlerTimeout)
	defer cancel()

	authResp, err := h.usecase.Execute(ctx, creds)
	if err != nil {
		slog.Error("login failed", "executionId", input.ExecutionId, "error", err)
		return models.LoginOutput{}, err
	}

	slog.Info("login succeeded", "executionId", input.ExecutionId)
	output := ToResponseAuth(authResp, input.ExecutionId)

	return output, nil
}

func ToResponseAuth(domainAuth domain.Auth, executionId string) models.LoginOutput {
	cookies := make([]http.Cookie, len(domainAuth.Cookies))
	for i, c := range domainAuth.Cookies {
		if c != nil {
			cookies[i] = *c
		}
	}
	return models.NewLoginOutput(cookies, domainAuth.InternalId, executionId)
}
