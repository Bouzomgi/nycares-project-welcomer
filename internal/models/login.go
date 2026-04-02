package models

import (
	"encoding/json"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

type Auth struct {
	Cookies []http.Cookie `json:"cookies"`
}

type LoginInput struct {
	ExecutionId string          `json:"executionId"`
	Context     json.RawMessage `json:"context,omitempty"`
}

type LoginOutput struct {
	Auth        Auth   `json:"auth"`
	InternalId  string `json:"internalId"`
	ExecutionId string `json:"executionId"`
}

func NewLoginOutput(cookies []http.Cookie, internalId string, executionId string) LoginOutput {
	return LoginOutput{
		Auth:        Auth{Cookies: cookies},
		InternalId:  internalId,
		ExecutionId: executionId,
	}
}

func ConvertAuth(a Auth) domain.Auth {
	out := domain.Auth{
		Cookies: make([]*http.Cookie, len(a.Cookies)),
	}

	for i := range a.Cookies {
		// Must take a copy inside the loop to avoid pointer aliasing issues
		c := a.Cookies[i]
		out.Cookies[i] = &c
	}

	return out
}
