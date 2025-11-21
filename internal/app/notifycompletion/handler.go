package notifycompletion

import (
	"context"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/domain"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

type NotifyCompletionHandler struct {
	usecase *NotifyCompletionUseCase
	cfg     *Config
}

func NewNotifyCompletionHandler(u *NotifyCompletionUseCase, cfg *Config) *NotifyCompletionHandler {
	return &NotifyCompletionHandler{usecase: u, cfg: cfg}
}

func (h *NotifyCompletionHandler) Handle(ctx context.Context, input models.NotifyCompletionInput) error {

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	domainProject, err := models.BuildDomainProject(input.Project)
	if err != nil {
		return err
	}

	notificationType, err := domain.ParseNotificationType(input.MessageToSend.Type)
	if err != nil {
		return err
	}

	err = h.usecase.Execute(ctx, notificationType, domainProject)
	if err != nil {
		return err
	}

	return nil
}
