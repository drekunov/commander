package mainwindow

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/widgets/config"
	"github.com/drekunov/gc/internal/widgets/info"
)

type Model struct {
	width, height int

	about *info.Model
}

func New() *Model {
	return &Model{
		about: info.New("Hello, World!", 0, 0),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	about, cmd := m.about.Update(msg)
	m.about, _ = about.(*info.Model)

	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
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
	exit := lipgloss.PlaceHorizontal(0, lipgloss.Center, "Press q to quit.")
	out := lipgloss.JoinVertical(lipgloss.Center, about, exit)
	out = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(out)

	return out
}
