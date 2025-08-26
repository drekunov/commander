package app

import (
	"fmt"

	"github.com/rivo/tview"
)

func Run() error {
	fmt.Println("Welcome to the g-commander. Orthodox file manager")

	app := tview.NewApplication()

	mainWindow := tview.NewForm()
	mainWindow.SetBorder(true).SetTitle("g-commander")

	mainWindow.AddButton("Quit", func() {
		app.Stop()
	})

	err := app.SetRoot(mainWindow, true).Run()
	if err != nil {
		return err
	}

	return nil
}
