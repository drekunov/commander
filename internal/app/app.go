package app

import (
	"fmt"

	"github.com/rivo/tview"
)

func Run() error {
	app := tview.NewApplication()

	mainWindow := tview.NewForm()
	mainWindow.SetBorder(true).SetTitle("g-commander")

	mainWindow.AddButton("Quit", func() {
		app.Stop()
	})

	err := app.SetRoot(mainWindow, true).Run()
	if err != nil {
		return fmt.Errorf("failed to run application: %w", err)
	}

	return nil
}
