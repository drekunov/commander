package app

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/drekunov/gc/internal/widgets/mainwindow"
)

func Run() error {
	_, err := tea.NewProgram(mainwindow.New(), tea.WithMouseAllMotion()).Run()
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return nil
}
