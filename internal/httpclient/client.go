package httpclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type Client struct {
	*http.Client
}

func NewClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	return &Client{Client: client}, nil
}

func (c *Client) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func (c *Client) ReadBody(resp *http.Response) ([]byte, error) {
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

func (c *Client) SetCookies(auth models.Auth, baseURL string) error {
	u, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %w", err)
	}

	c.Jar.SetCookies(u, auth.Cookies)
	return nil
}

func (c *Client) GetCookies(baseURL string) ([]*http.Cookie, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	return c.Jar.Cookies(u), nil
}
