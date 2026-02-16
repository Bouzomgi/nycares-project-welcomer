package endpoints

import "os"

var BaseUrl = "https://www.newyorkcares.org"

func init() {
	if url := os.Getenv("NYCARES_API_BASE_URL"); url != "" {
		BaseUrl = url
	}
}

const (
	LoginPath       = "/user/login"
	GetSchedulePath = "/api/schedule/retrieve"
	GetCampaignPath = "/api/campaign/retrieve"
)
