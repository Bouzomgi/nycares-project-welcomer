package httpservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

type HttpService struct {
	client *http.Client
}

type HttpClientOption func(*http.Client) error

func NewHttpService(opts ...HttpClientOption) (*HttpService, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	return &HttpService{client}, nil
}

func WithAuth(auth domain.Auth, baseURL string) HttpClientOption {
	return func(client *http.Client) error {
		u, err := url.Parse(baseURL)
		if err != nil {
			return fmt.Errorf("failed to parse base URL: %w", err)
		}

		client.Jar.SetCookies(u, auth.Cookies)
		return nil
	}
}

func (s *HttpService) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func (s *HttpService) ReadBody(resp *http.Response) ([]byte, error) {
	if resp.Body == nil {
		return nil, fmt.Errorf("response body is nil")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func (s *HttpService) GetCookies(baseURL string) ([]*http.Cookie, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	return s.client.Jar.Cookies(u), nil
}
