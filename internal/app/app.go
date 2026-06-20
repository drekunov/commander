package app

import (
	"context"
	"fmt"
	"log"
)

type App struct {
	ui UI

	connector Connector
}

func New(ui UI, conn Connector) *App {
	return &App{
		ui:        ui,
		connector: conn,
	}
}

func (a *App) Run(ctx context.Context) error {
	log.Println(getBuildInfo())

	data, err := a.connector.ReadDir("/")
	if err != nil {
		a.ui.Info(ctx, a.connector.Name(), "/", err.Error())

		return fmt.Errorf("error reading directory: %w", err)
	}

	a.ui.SetData(data)

	<-ctx.Done()

	return nil
}
