package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	uiInstance := ui.New()
	appInstance := app.New(uiInstance)

	errGr, ctx := errgroup.WithContext(ctx)

	errGr.Go(func() error {
		err := uiInstance.Run(ctx)
		if err != nil {
			return fmt.Errorf("could not start ui: %w", err)
		}

		return nil
	})

	errGr.Go(func() error {
		err := appInstance.Run(ctx)
		if err != nil {
			return err
		}

		return nil
	})

	err := errGr.Wait()
	if err != nil {
		return err
	}

	return nil
}
