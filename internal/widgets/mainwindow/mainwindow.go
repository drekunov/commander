package mainwindow

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/widgets/config"
	"github.com/drekunov/gc/internal/widgets/info"
)

type Model struct {
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

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	about, cmd := m.about.Update(msg)
	m.about, _ = about.(*info.Model)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	about := m.about.View()
	exit := lipgloss.PlaceHorizontal(0, lipgloss.Center, "Press q to quit.")
	out := lipgloss.JoinVertical(lipgloss.Center, about, exit)
	out = config.Values.DialogBoxStyle.Render(out)

	return out
}
