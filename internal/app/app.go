package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/drekunov/gc/internal/widgets/mainwindow"
)

func Run() error {
	if _, err := tea.NewProgram(mainwindow.New(), tea.WithMouseAllMotion()).Run(); err != nil {
		return err
	}

	return nil
}
