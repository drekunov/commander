package mainwindow

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

type Model struct {
	width, height int

	windows windowsList
}

func New() *Model {
	return &Model{
		windows: make(windowsList),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.windows, cmd = m.windows.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 2
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	windows := m.windows.View()

	exit := lipgloss.PlaceHorizontal(0, lipgloss.Center, "Press F10 to quit.")
	out := lipgloss.JoinVertical(lipgloss.Center, windows, exit)
	out = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(out)

	return out
}
