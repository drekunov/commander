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

	go func() {
		a.ui.Info(ctx, "Hello", "footer", msg)
	}()

	go func() {
		msg = a.ui.Input(ctx, "test header", "test footer", msg)
	}()

	<-ctx.Done()

	return nil
}
