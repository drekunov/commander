package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/connectors/filesystem"
	"github.com/drekunov/gc/internal/ui"
	"golang.org/x/sync/errgroup"
)

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}

func Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fsConn := filesystem.New()

	uiInstance := ui.New()
	appInstance := app.New(uiInstance, fsConn)

	errGr, _ := errgroup.WithContext(context.Background())

	errGr.Go(func() error {
		err := uiInstance.Run()
		cancel()
		if err != nil && !errors.Is(err, tea.ErrInterrupted) {
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
	if err != nil {
		return fmt.Errorf("exiting app: %w", err)
	}

	return nil
}
