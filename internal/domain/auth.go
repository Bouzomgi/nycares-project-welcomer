package domain

import "net/http"

type Credentials struct {
	Username string
	Password string
}

type Auth struct {
	Cookies []*http.Cookie
}
