package models

import "net/http"

type Auth struct {
	Cookies []http.Cookie `json:"cookies"`
}
