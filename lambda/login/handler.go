package main

import (
	"context"
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

func (h *LoginHandler) Handle(ctx context.Context) (models.LoginOutput, error) {
	creds := domain.Credentials{
		Username: h.cfg.Account.Username,
		Password: h.cfg.Account.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, config.HTTPHandlerTimeout)
	defer cancel()

	authResp, err := h.usecase.Execute(ctx, creds)
	if err != nil {
		return models.LoginOutput{}, err
	}

	output := ToResponseAuth(authResp)

	return output, nil
}

func ToResponseAuth(domainAuth domain.Auth) models.LoginOutput {
	cookies := make([]http.Cookie, len(domainAuth.Cookies))
	for i, c := range domainAuth.Cookies {
		if c != nil {
			cookies[i] = *c
		}
	}
	return models.NewLoginOutput(cookies)
}
