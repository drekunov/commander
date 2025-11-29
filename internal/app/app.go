package app

import (
	"context"
	"time"
)

type App struct {
	ui Dialogs
}

func New(ui Dialogs) *App {
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
			time.Sleep(time.Second)
			a.ui.Info("Hello", "World")
		}
	}
}
