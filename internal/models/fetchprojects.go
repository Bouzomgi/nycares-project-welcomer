package models

import (
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
)

type FetchProjectsOutput struct {
	Auth     authData      `json:"auth"`
	Projects []projectData `json:"projects"`
}

type authData struct {
	Cookies []http.Cookie `json:"cookies"`
}

type projectData struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

func BuildFetchProjectsOutput(input FetchProjectsInput, domainProjects []domain.Project) FetchProjectsOutput {
	projects := make([]projectData, len(domainProjects))
	for i, p := range domainProjects {
		projects[i] = projectData{
			Name: p.Name,
			Date: p.Date,
		}
	}

	return FetchProjectsOutput{
		Auth: authData{
			Cookies: input.Auth.Cookies,
		},
		Projects: projects,
	}
}
