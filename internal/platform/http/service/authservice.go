package httpservice

import (
	"bytes"
	"context"
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

	resp, err := s.SendRequest(ctx, req)
	if err != nil {
		return domain.Auth{}, fmt.Errorf("login request failed: %w", err)
	}

	if err := CheckResponse(resp); err != nil {
		return domain.Auth{}, fmt.Errorf("login failed: %w", err)
	}

	body, err := s.ReadBody(resp)
	if err != nil {
		return domain.Auth{}, fmt.Errorf("failed to read login response body: %w", err)
	}

	re := regexp.MustCompile(`schedule/retrieve/([A-Za-z0-9]{15,18})`)
	matches := re.FindSubmatch(body)
	if matches == nil {
		bodyLen := len(body)
		snippet := body
		if bodyLen > 500 {
			snippet = body[:500]
		}
		slog.Error("internalId not found in login response",
			"finalUrl", resp.Request.URL.String(),
			"contentType", resp.Header.Get("Content-Type"),
			"bodyLen", bodyLen,
			"bodySnippet", string(snippet),
		)
		return domain.Auth{}, fmt.Errorf("internalId not found in login response")
	}
	internalId := string(matches[1])

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
	writer.Close()

	loginURL := endpoints.JoinPaths(s.baseUrl, endpoints.LoginPath)

	req, err := http.NewRequest("POST", loginURL, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}
