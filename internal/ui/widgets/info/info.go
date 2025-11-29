package info

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/ui/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Model struct {
	title   string
	text    string
	width   int
	height  int
	visible bool

	okButton *button.Model
}

func New() *Model {
	return &Model{
		okButton: button.New("Ok", true),
	}
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) SetText(text string) {
	m.text = text
}

func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.okButton.Update(msg)

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

	okButton := m.okButton.View()

	output := lipgloss.JoinVertical(lipgloss.Center, text, okButton)

	output = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(output)

	output = lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		output,
	)

	return output
}
