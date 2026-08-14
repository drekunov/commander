package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

// Confirm presents a yes/no choice and reports whether the user accepted.
type Confirm struct {
	dialogBase

	message string
	cursor  int

	decision bool
	done     chan struct{}
}

func NewConfirm(styles config.Styles, message string) *Confirm {
	return &Confirm{
		dialogBase: newDialogBase(styles),
		message:    message,
		done:       make(chan struct{}, 1),
	}
}

func (m *Confirm) Init() tea.Cmd {
	return nil
}

func (m *Confirm) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft, tea.KeyRight:
			m.cursor = 1 - m.cursor
		case tea.KeyEnter:
			if m.visible {
				m.decision = m.cursor == 0
				m.visible = false

				m.done <- struct{}{}
			}
		}
	}

	return m, nil
}

func (m *Confirm) View() string {
	if !m.visible {
		return ""
	}

	message := m.styles.TextStyle.Render(m.message)

	buttons := lipgloss.JoinHorizontal(
		lipgloss.Center,
		m.buttonView("Yes", m.cursor == 0),
		m.buttonView("No", m.cursor == 1),
	)

	output := lipgloss.JoinVertical(lipgloss.Center, message, buttons)

	return m.render(output)
}

func (m *Confirm) Decision() bool {
	return m.decision
}

func (m *Confirm) Done() chan struct{} {
	return m.done
}

func (m *Confirm) buttonView(label string, focused bool) string {
	style := m.styles.ButtonStyle
	if focused {
		style = m.styles.ActiveButtonStyle
	}

	return style.Render(label)
}
