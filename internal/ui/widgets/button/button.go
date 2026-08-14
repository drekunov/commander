package button

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
)

type Model struct {
	text    string
	pressed bool
	focused bool

	styles config.Styles
}

func New(text string, focused bool, styles config.Styles) *Model {
	return &Model{
		text:    text,
		focused: focused,
		styles:  styles,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
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
	button := m.styles.ButtonStyle.
		Render(m.text)

	if m.focused {
		button = m.styles.ActiveButtonStyle.Render(m.text)
	}

	if m.pressed {
		button = m.styles.PressedButtonStyle.Render(m.text)
		m.pressed = false
	}

	return button
}
