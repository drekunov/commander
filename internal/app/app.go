package app

import (
	"context"
)

type App struct {
	ui UI
}

func New(ui UI) *App {
	return &App{
		ui: ui,
	}
}

func (a *App) Run(ctx context.Context) error {
	msg := "Hello World!!!"

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg = a.ui.Input(ctx, "test header", "test footer", msg)
		}
	}
}
