package fetchprojects

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	httpservice "github.com/Bouzomgi/nycares-project-welcomer/internal/platform/http/service"
)

// AuthFailureException is returned when the upcoming projects request gets a 401/403,
// indicating expired auth cookies. Step Functions catches this by type name
// to route back to Login for a fresh authentication attempt.
type AuthFailureException struct{}

func (e *AuthFailureException) Error() string { return "auth failure: cookies may be expired" }

type FetchProjectsUseCase struct {
	upcomingSvc httpservice.UpcomingProjectsService
}

func NewFetchProjectsUseCase(upcomingSvc httpservice.UpcomingProjectsService) *FetchProjectsUseCase {
	return &FetchProjectsUseCase{upcomingSvc: upcomingSvc}
}

func (u *FetchProjectsUseCase) Execute(ctx context.Context, auth domain.Auth, userSFID string) ([]domain.Project, error) {

	u.upcomingSvc.SetCookies(auth.Cookies)
	slog.Info("fetchprojects invoking upcoming projects", "cookieCount", len(auth.Cookies))
	projects, err := u.upcomingSvc.GetUpcomingProjects(ctx, userSFID)
	if err != nil {
		var httpErr *httpservice.HTTPError
		if errors.As(err, &httpErr) && (httpErr.StatusCode == http.StatusUnauthorized || httpErr.StatusCode == http.StatusForbidden) {
			slog.Error("upcoming projects request auth failure", "statusCode", httpErr.StatusCode, "url", httpErr.URL, "body", string(httpErr.Body))
			return nil, &AuthFailureException{}
		}
		return nil, err
	}

	return projects, nil
}
