package api

// API endpoints and paths
const (
	// BaseURL is the root URL for the API
	BaseURL = "https://www.newyorkcares.org"

	// Auth endpoints
	LoginPath = "/user/login"

	// Project endpoints
	GetSchedulePath = "/api/schedule/retrieve"
)

func BuildURL(baseURL, path string) string {
	return baseURL + path
}
