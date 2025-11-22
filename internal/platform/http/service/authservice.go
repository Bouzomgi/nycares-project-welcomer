package httpservice

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/endpoints"
)

type AuthService interface {
	Login(ctx context.Context, creds domain.Credentials) (domain.Auth, error)
}

func (s *HttpService) Login(ctx context.Context, creds domain.Credentials) (domain.Auth, error) {
	if creds.Username == "" {
		return domain.Auth{}, fmt.Errorf("username is required")
	}

	if creds.Password == "" {
		return domain.Auth{}, fmt.Errorf("password is required")
	}

	req, err := s.buildLoginRequest(creds)
	if err != nil {
		return domain.Auth{}, fmt.Errorf("failed to build login request: %w", err)
	}

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return domain.Auth{}, fmt.Errorf("login request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return domain.Auth{}, fmt.Errorf("login failed: %w", err)
	}

	cookies, err := s.GetCookies()
	if err != nil {
		return domain.Auth{}, fmt.Errorf("failed to get cookies: %w", err)
	}

	if len(cookies) == 0 {
		return domain.Auth{}, fmt.Errorf("no cookies set after login")
	}

	return domain.Auth{Cookies: cookies}, nil
}

func (s *HttpService) buildLoginRequest(creds domain.Credentials) (*http.Request, error) {
	form := url.Values{}
	form.Set("form_id", "user_login_form")
	form.Set("name", creds.Username)
	form.Set("pass", creds.Password)

	encoded := form.Encode()
	loginURL := endpoints.JoinPaths(s.baseUrl, endpoints.LoginPath)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}
