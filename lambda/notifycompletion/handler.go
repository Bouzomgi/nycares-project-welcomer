package main

import (
	"context"
	"log/slog"

	nc "github.com/Bouzomgi/nycares-project-welcomer/internal/app/notifycompletion"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type NotifyCompletionHandler struct {
	usecase *nc.NotifyCompletionUseCase
	cfg     *nc.Config
}

func NewNotifyCompletionHandler(u *nc.NotifyCompletionUseCase, cfg *nc.Config) *NotifyCompletionHandler {
	return &NotifyCompletionHandler{usecase: u, cfg: cfg}
}

func (h *NotifyCompletionHandler) Handle(ctx context.Context, input models.NotifyCompletionInput) error {
	slog.Info("notifycompletion handler invoked")

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	domainProject, err := models.ConvertProjectNotificationToDomainProject(input.ExistingProjectNotification)
	if err != nil {
		slog.Error("notifycompletion failed to convert project", "error", err)
		return err
	}

	notificationType, err := domain.ParseNotificationType(input.MessageToSend.Type)
	if err != nil {
		slog.Error("notifycompletion invalid notification type", "error", err)
		return err
	}

	err = h.usecase.Execute(ctx, notificationType, domainProject)
	if err != nil {
		slog.Error("notifycompletion failed", "error", err)
		return err
	}

	slog.Info("notifycompletion succeeded")
	return nil
}
