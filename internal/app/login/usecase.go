package login

import (
	"context"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
)

type LoginUseCase struct {
	authSvc httpservice.AuthService
}

// TODO: add specific services like below everywhere
func NewLoginUseCase(authSvc httpservice.AuthService) *LoginUseCase {
	return &LoginUseCase{authSvc: authSvc}
}

func (u *LoginUseCase) Execute(ctx context.Context, creds domain.Credentials) (domain.Auth, error) {

	auth, err := u.authSvc.Login(ctx, creds)
	if err != nil {
		return domain.Auth{}, err
	}

	return domain.Auth{
		Cookies: auth.Cookies,
	}, nil
}
