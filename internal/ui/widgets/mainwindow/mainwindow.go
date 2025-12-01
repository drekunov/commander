package mainwindow

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/info"
)

type Model struct {
	width, height int

	about *info.Model
}

func New() *Model {
	return &Model{
		about: info.New(),
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

	m.about, cmd = m.about.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			m.about.Free()

			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width - 2
		m.height = msg.Height - 2
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	about := m.about.View()
	exit := lipgloss.PlaceHorizontal(0, lipgloss.Center, "Press F10 to quit.")
	out := lipgloss.JoinVertical(lipgloss.Center, about, exit)
	out = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(out)

	return out
}
