package httpservice

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/httpclient"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/service/api"
)

type AuthService interface {
	Login(ctx context.Context, creds models.Credentials) (models.Auth, error)
}

func WithAuth(auth models.Auth) HttpServiceOption {
	return func(s *HttpService) error {
		return s.client.SetCookies(auth, s.baseURL)
	}
}

func (s *HttpService) Login(ctx context.Context, creds models.Credentials) (models.Auth, error) {
	if creds.Username == "" {
		return models.Auth{}, fmt.Errorf("username is required")
	}

	if creds.Password == "" {
		return models.Auth{}, fmt.Errorf("password is required")
	}

	req, err := s.buildLoginRequest(creds)
	if err != nil {
		return models.Auth{}, fmt.Errorf("failed to build login request: %w", err)
	}

	resp, err := s.client.SendRequest(ctx, req)
	if err != nil {
		return models.Auth{}, fmt.Errorf("login request failed: %w", err)
	}

	if err := httpclient.CheckResponse(resp); err != nil {
		return models.Auth{}, fmt.Errorf("login failed: %w", err)
	}

	cookies, err := s.client.GetCookies(s.baseURL)
	if err != nil {
		return models.Auth{}, fmt.Errorf("failed to get cookies: %w", err)
	}

	if len(cookies) == 0 {
		return models.Auth{}, fmt.Errorf("no cookies set after login")
	}

	return s.createAuthFromCookies(cookies), nil
}

func (s *HttpService) buildLoginRequest(creds models.Credentials) (*http.Request, error) {
	form := url.Values{}
	form.Set("form_id", "user_login_form")
	form.Set("name", creds.Username)
	form.Set("pass", creds.Password)

	encoded := form.Encode()
	loginURL := api.BuildURL(s.baseURL, api.LoginPath)

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(encoded))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

func (s *HttpService) extractCookies() ([]*http.Cookie, error) {
	u, err := url.Parse(s.baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	cookies := s.client.Jar.Cookies(u)
	if len(cookies) == 0 {
		return nil, fmt.Errorf("login failed: no cookies set on %s", s.baseURL)
	}

	return cookies, nil
}

func (s *HttpService) createAuthFromCookies(cookies []*http.Cookie) models.Auth {
	var auth models.Auth
	for _, c := range cookies {
		auth.Cookies = append(auth.Cookies, &http.Cookie{
			Name:   c.Name,
			Value:  c.Value,
			Domain: c.Domain,
			Path:   c.Path,
		})
	}
	return auth
}
