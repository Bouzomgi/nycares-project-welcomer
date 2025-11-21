package httpservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
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
		Jar: jar,
	}

	return &HttpService{client, baseUrl}, nil
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

func (s *HttpService) GetCookies() ([]*http.Cookie, error) {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	return s.client.Jar.Cookies(u), nil // nukes the expr date!!
}

func (s *HttpService) SetCookies(cookies []*http.Cookie) error {
	u, err := url.Parse(s.baseUrl)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	s.client.Jar.SetCookies(u, cookies)
	return nil
}
