package domain

import (
	"encoding/json"
	"net/http"
)

type Credentials struct {
	Username         string
	Password         string
	MockProjectsJSON json.RawMessage
}

type Auth struct {
	Cookies    []*http.Cookie
	InternalId string
}
