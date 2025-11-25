package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/widgets/config"
)

type Model struct {
	text    string
	pressed bool
	focused bool
}

func New(text string, focused bool) *Model {
	return &Model{
		text:    text,
		focused: focused,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.pressed = true
		}
	}

	return m, nil
}

func (m *Model) View() string {
	button := config.Values.ButtonStyle.
		Render(m.text)

	if m.focused {
		button = config.Values.ActiveButtonStyle.Render(m.text)
	}

	if m.pressed {
		button = config.Values.PressedButtonStyle.Render(m.text)
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.BlockBorder(), false, true, true, false).
		BorderForeground(lipgloss.Color("#111111")).
		BorderBackground(lipgloss.Color("#111111"))

	return box.Render(button)
}
