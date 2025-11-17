package httpservice

import (
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/httpclient"
)

type HttpService struct {
	client  *httpclient.Client
	baseURL string
}

type HttpServiceOption func(*HttpService) error

func NewHttpService(baseURL string, opts ...HttpServiceOption) (*HttpService, error) {
	client, err := httpclient.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create http client: %w", err)
	}
	service := &HttpService{
		client: client, baseURL: baseURL,
	}
	for _, opt := range opts {
		if err := opt(service); err != nil {
			return nil, err
		}
	}
	return service, nil
}
