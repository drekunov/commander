package info

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/widgets/config"
)

type Model struct {
	text    string
	width   int
	height  int
	visible bool
}

func New(text string, width, height int) *Model {
	return &Model{
		text:    text,
		width:   width,
		height:  height,
		visible: true,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
			m.visible = false
		}
	}

	return m, nil
}

func (m *Model) View() string {
	if !m.visible {
		return ""
	}

	text := config.Values.TextStyle.Render(m.text)
	okButton := config.Values.ActiveButtonStyle.Render("Ok")

	output := lipgloss.JoinVertical(lipgloss.Center, text, okButton)

	output = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(output)

	dialog := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		output,
	)

	return dialog
}
