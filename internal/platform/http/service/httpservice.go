package httpservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type HttpService struct {
	client  *http.Client
	baseUrl string
}

func NewHttpService(baseUrl string) (*HttpService, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &HttpService{client, baseUrl}, nil
}

func (s *HttpService) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	const maxAttempts = 3
	baseDelay := time.Second

	req = req.WithContext(ctx)

	var lastErr error
	for attempt := range maxAttempts {
		if attempt > 0 {
			if req.GetBody != nil {
				body, err := req.GetBody()
				if err != nil {
					return nil, fmt.Errorf("failed to reset request body for retry: %w", err)
				}
				req.Body = body
			}
			delay := baseDelay * time.Duration(1<<(attempt-1))
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			continue
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("request failed with status %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, lastErr
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

func (s *HttpService) GetCookies() ([]*http.Cookie, error) {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	return s.client.Jar.Cookies(u), nil
}

func (s *HttpService) SetCookies(cookies []*http.Cookie) error {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	s.client.Jar.SetCookies(u, cookies)
	return nil
}
