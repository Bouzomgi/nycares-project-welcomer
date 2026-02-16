package main

import (
	"context"
	"log/slog"

	dlq "github.com/Bouzomgi/nycares-project-welcomer/internal/app/dlqnotifier"
	"github.com/Bouzomgi/nycares-project-welcomer/internal/config"
)

type DLQNotifierHandler struct {
	usecase *dlq.DLQNotifierUseCase
	cfg     *dlq.Config
}

func NewDLQNotifierHandler(u *dlq.DLQNotifierUseCase, cfg *dlq.Config) *DLQNotifierHandler {
	return &DLQNotifierHandler{usecase: u, cfg: cfg}
}

func (h *DLQNotifierHandler) Handle(ctx context.Context, input map[string]interface{}) error {
	slog.Error("dlqnotifier handler invoked", "input", input)

	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	err := h.usecase.Execute(ctx, input)
	if err != nil {
		slog.Error("dlqnotifier failed to send notification", "error", err)
	}
	return err
}
