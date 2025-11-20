package login

import (
	"context"
	"net/http"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type LoginHandler struct {
	usecase *LoginUseCase
	cfg     *Config
}

func NewLoginHandler(u *LoginUseCase, cfg *Config) *LoginHandler {
	return &LoginHandler{usecase: u, cfg: cfg}
}

func (h *LoginHandler) Handle(ctx context.Context) (models.LoginOutput, error) {
	creds := domain.Credentials{
		Username: h.cfg.Account.Username,
		Password: h.cfg.Account.Password,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
