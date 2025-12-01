package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/mainwindow"
)

type Model struct {
	program *tea.Program

	main *mainwindow.Model
}

func New() *Model {
	model := &Model{
		main: mainwindow.New(),
	}

	model.program = tea.NewProgram(model, tea.WithMouseAllMotion(), tea.WithAltScreen())

	return model
}

func (m *Model) Run() error {
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
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.main, cmd = m.main.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	return m.main.View()
}
