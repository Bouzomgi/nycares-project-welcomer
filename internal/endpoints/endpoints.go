package endpoints

const BaseUrl = "https://www.newyorkcares.org"

const (
	LoginPath               = "/user/login"
	GetSchedulePath         = "/api/schedule/retrieve" // deprecated: use GetUpcomingProjectsPath
	GetUpcomingProjectsPath = "/api/registrations/dashboard/upcoming"
	GetCampaignPath         = "/api/campaign/retrieve"
)
