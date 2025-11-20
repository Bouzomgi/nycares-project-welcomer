package fetchprojects

import (
	"context"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http"
)

type FetchProjectsUseCase struct {
	scheduleSvc httpservice.ScheduleService
}

func NewFetchProjectsUseCase(scheduleSvc httpservice.ScheduleService) *FetchProjectsUseCase {
	return &FetchProjectsUseCase{scheduleSvc: scheduleSvc}
}

func (u *FetchProjectsUseCase) Execute(ctx context.Context, auth domain.Auth, internalID string) ([]domain.Project, error) {

	u.scheduleSvc.SetCookies(auth.Cookies)
	projects, err := u.scheduleSvc.GetSchedule(ctx, internalID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}
