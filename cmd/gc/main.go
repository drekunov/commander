package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/ui"
	"golang.org/x/sync/errgroup"
)

func main() {
	err := Run()
	if err != nil {
		log.Print(err)
	}
}

func Run() error {
	ctx := context.Background()

	uiInstance := ui.New()
	appInstance := app.New(uiInstance)

	errGr, ctx := errgroup.WithContext(ctx)

	errGr.Go(func() error {
		err := uiInstance.Run()
		if err != nil {
			return fmt.Errorf("could not start ui: %w", err)
		}

		return nil
	})

	errGr.Go(func() error {
		err := appInstance.Run(ctx)
		if err != nil {
			return fmt.Errorf("could not start app: %w", err)
		}

		return nil
	})

	err := errGr.Wait()
	if err != nil && !errors.Is(err, tea.ErrInterrupted) {
		return fmt.Errorf("exiting app: %w", err)
	}

	return nil
}
