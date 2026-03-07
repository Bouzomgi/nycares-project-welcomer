package fetchprojects

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
)

// AuthFailureException is returned when the schedule request gets a 401/403,
// indicating expired auth cookies. Step Functions catches this by type name
// to route back to Login for a fresh authentication attempt.
type AuthFailureException struct{}

func (e *AuthFailureException) Error() string { return "auth failure: cookies may be expired" }

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
		var httpErr *httpservice.HTTPError
		if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusUnauthorized || httpErr.StatusCode == http.StatusForbidden) {
			return nil, &AuthFailureException{}
		}
		return nil, err
	}

	return projects, nil
}
