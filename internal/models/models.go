package models

type Cookie struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Domain string `json:"domain"`
	Path   string `json:"path"`
}

type Auth struct {
	Cookies []Cookie `json:"cookies"`
}
