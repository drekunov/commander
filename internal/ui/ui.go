package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

type Model struct {
	program *tea.Program

	main *mainform.Model
	wm   *wm.Manager
}

func New() *Model {
	model := &Model{
		main: mainform.New(),
		wm:   wm.New(),
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
	return m.main.Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0, 2)

	var cmd tea.Cmd

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		width := msg.Width - 2
		height := msg.Height - 2
		m.wm.SetSize(width, height)
	}

	// Update window manager (handles mouse, keys to focused window).
	cmd = m.wm.Update(msg)
	cmds = append(cmds, cmd)

	// Update main form (panel, etc.) only if no windows have focus.
	m.main, cmd = m.main.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	base := m.main.View()

	wmView := m.wm.View()
	if wmView == "" {
		return base
	}

	// Overlay wm windows on top of main form via canvas compositing.
	width, height := m.wm.Width(), m.wm.Height()

	return wm.Overlay(base, wmView, width, height)
}
