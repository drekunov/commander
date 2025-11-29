package ui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/mainwindow"
)

type Model struct {
	program *tea.Program

	main *mainwindow.Model
}

func New() *Model {
	return &Model{
		main: mainwindow.New(),
	}
}

func (m *Model) Run(ctx context.Context) error {
	m.program = tea.NewProgram(m, tea.WithContext(ctx), tea.WithMouseAllMotion())

	_, err := m.program.Run()
	if err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}

	return tea.ErrInterrupted
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.main.Update(msg)
}

func (m *Model) View() string {
	return m.main.View()
}
