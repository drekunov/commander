package app

import (
	"context"
	"time"
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
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(100 * time.Millisecond)
			a.ui.Info("test header", "test footer", "Hello World!!!")
		}
	}
}
