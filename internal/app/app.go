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
	a.ui.Input("test header", "test footer", "Hello World!!!")

	<-ctx.Done()

	return nil
}
