package models

import (
	"fmt"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/utils"
)

type FetchProjectsInput = LoginOutput

type FetchProjectsOutput struct {
	Auth     Auth      `json:"auth"`
	Projects []project `json:"projects"`
}

type project struct {
	Name string `json:"name"`
	Date string `json:"date"`
	Id   string `json:"id"`
}

// MODEL -> DOMAIN
func BuildDomainProject(p project) (domain.Project, error) {
	projectDate, err := utils.StringToDate(p.Date)

	if err != nil {
		return domain.Project{}, fmt.Errorf("could not parse project date")
	}

	domainProject := domain.Project{
		Name: p.Name,
		Date: projectDate,
		Id:   p.Id,
	}

	return domainProject, nil
}

// DOMAIN -> MODEL
func buildModelProject(p domain.Project) project {
	return project{
		Name: p.Name,
		Date: utils.DateToString(p.Date),
		Id:   p.Id,
	}
}

func BuildFetchProjectsOutput(input FetchProjectsInput, domainProjects []domain.Project) FetchProjectsOutput {
	projects := make([]project, len(domainProjects))
	for i, p := range domainProjects {
		projects[i] = buildModelProject(p)
	}

	return FetchProjectsOutput{
		Auth: Auth{
			Cookies: input.Auth.Cookies,
		},
		Projects: projects,
	}
}
