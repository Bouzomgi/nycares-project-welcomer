package main

import (
	"context"
	"log/slog"

	fp "github.com/Bouzomgi/nycares-project-welcomer/internal/app/fetchprojects"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type FetchProjectsHandler struct {
	usecase *fp.FetchProjectsUseCase
	cfg     *fp.Config
}

func NewFetchProjectsHandler(u *fp.FetchProjectsUseCase, cfg *fp.Config) *FetchProjectsHandler {
	return &FetchProjectsHandler{usecase: u, cfg: cfg}
}

func (h *FetchProjectsHandler) Handle(ctx context.Context, input models.FetchProjectsInput) (models.FetchProjectsOutput, error) {
	slog.Info("fetchprojects handler invoked")

	ctx, cancel := context.WithTimeout(ctx, config.HTTPHandlerTimeout)
	defer cancel()

	auth := models.ConvertAuth(input.Auth)

	projects, err := h.usecase.Execute(ctx, auth, h.cfg.Account.InternalId)
	if err != nil {
		slog.Error("fetchprojects failed", "error", err)
		return models.FetchProjectsOutput{}, err
	}

	slog.Info("fetchprojects succeeded", "count", len(projects))
	output := models.BuildFetchProjectsOutput(input, projects)

	return output, nil
}
