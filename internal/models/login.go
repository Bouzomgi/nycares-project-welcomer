package models

import (
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

type Auth struct {
	Cookies []http.Cookie `json:"cookies"`
}

type LoginOutput struct {
	Auth Auth `json:"auth"`
}

func NewLoginOutput(cookies []http.Cookie) LoginOutput {
	return LoginOutput{
		Auth{
			Cookies: cookies,
		},
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
