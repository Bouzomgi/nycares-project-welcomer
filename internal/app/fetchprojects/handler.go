package fetchprojects

import (
	"context"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type FetchProjectsHandler struct {
	usecase *FetchProjectsUseCase
	cfg     *Config
}

func NewFetchProjectsHandler(u *FetchProjectsUseCase, cfg *Config) *FetchProjectsHandler {
	return &FetchProjectsHandler{usecase: u, cfg: cfg}
}

func (h *FetchProjectsHandler) Handle(ctx context.Context, input models.FetchProjectsInput) (models.FetchProjectsOutput, error) {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	auth := models.ConvertAuth(input.Auth)

	projects, err := h.usecase.Execute(ctx, auth, h.cfg.Account.InternalId)
	if err != nil {
		return models.FetchProjectsOutput{}, err
	}

	output := models.BuildFetchProjectsOutput(input, projects)

	return output, nil
}
