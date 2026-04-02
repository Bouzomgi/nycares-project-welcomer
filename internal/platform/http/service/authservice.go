package httpservice

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"regexp"

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

	var redirectURLs []string
	s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectURLs = append(redirectURLs, req.URL.String())
		return nil
	}
	defer func() { s.client.CheckRedirect = nil }()

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return domain.Auth{}, fmt.Errorf("login request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return domain.Auth{}, fmt.Errorf("login failed: %w", err)
	}

	re := regexp.MustCompile(`003[A-Za-z0-9]{15}`)

	var internalId string
	for _, u := range redirectURLs {
		if m := re.FindString(u); m != "" {
			internalId = m
			break
		}
	}

	if internalId == "" {
		body, err := s.ReadBody(resp)
		if err != nil {
			return domain.Auth{}, fmt.Errorf("failed to read login response body: %w", err)
		}
		if m := re.Find(body); m != nil {
			internalId = string(m)
		}
	}

	if internalId == "" {
		slog.Error("internalId not found in login response", "redirectUrls", redirectURLs, "finalUrl", resp.Request.URL.String())
		return domain.Auth{}, fmt.Errorf("internalId not found in login response")
	}

	cookies, err := s.GetCookies()
	if err != nil {
		return domain.Auth{}, fmt.Errorf("failed to get cookies: %w", err)
	}

	if len(cookies) == 0 {
		return domain.Auth{}, fmt.Errorf("no cookies set after login")
	}

	return domain.Auth{Cookies: cookies, InternalId: internalId}, nil
}

func (s *HttpService) buildLoginRequest(creds domain.Credentials) (*http.Request, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	writer.WriteField("form_id", "user_login_form")
	writer.WriteField("name", creds.Username)
	writer.WriteField("pass", creds.Password)
	// mock_projects is read by the mock server to return per-execution project config
	// via the session cookie. Empty in production — the real API ignores unknown fields.
	if len(creds.MockProjectsJSON) > 0 {
		writer.WriteField("mock_projects", base64.StdEncoding.EncodeToString(creds.MockProjectsJSON))
	}
	writer.Close()

	loginURL := endpoints.JoinPaths(s.baseUrl, endpoints.LoginPath)

	req, err := http.NewRequest("POST", loginURL, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
