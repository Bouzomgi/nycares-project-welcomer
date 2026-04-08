package main

import (
	"context"
	"log/slog"
	"time"

	rp "github.com/Bouzomgi/nycares-project-welcomer/internal/app/routeproject"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type RouteProjectHandler struct {
	usecase *rp.RouteProjectUseCase
	cfg     *rp.Config
}

func NewRouteProjectHandler(u *rp.RouteProjectUseCase, cfg *rp.Config) *RouteProjectHandler {
	return &RouteProjectHandler{usecase: u, cfg: cfg}
}

func (h *RouteProjectHandler) Handle(ctx context.Context, input models.RouteProjectInput) (models.RouteProjectOutput, error) {
	slog.Info("routeproject handler invoked", "executionId", input.ExecutionId)

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	domainProject, err := models.BuildDomainProject(input.Project)
	if err != nil {
		slog.Error("routeproject failed to build project", "executionId", input.ExecutionId, "error", err)
		return models.RouteProjectOutput{}, err
	}

	existingNotification, messageType, targetSendTime, err := h.usecase.Execute(ctx, domainProject)
	if err != nil {
		slog.Error("routeproject failed", "executionId", input.ExecutionId, "error", err)
		return models.RouteProjectOutput{}, err
	}

	slog.Info("routeproject succeeded", "executionId", input.ExecutionId, "messageType", messageType.String())

	outputNotification := models.ConvertDomainProjectNotification(existingNotification)
	outputNotification.ChannelId = input.Project.ChannelId

	targetSendTimeStr := ""
	if !targetSendTime.IsZero() {
		targetSendTimeStr = targetSendTime.UTC().Format(time.RFC3339)
	}

	return models.RouteProjectOutput{
		Auth:                        input.Auth,
		ExistingProjectNotification: outputNotification,
		MessageType:                 messageType.String(),
		TargetSendTime:              targetSendTimeStr,
		ExecutionId:                 input.ExecutionId,
	}, nil
}
