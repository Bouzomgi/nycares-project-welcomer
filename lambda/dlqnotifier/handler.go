package main

import (
	"context"

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
	ctx, cancel := context.WithTimeout(ctx, config.DefaultHandlerTimeout)
	defer cancel()

	return h.usecase.Execute(ctx, input)
}
