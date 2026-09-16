package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Info struct {
	dialogBase

	text string

	okButton *button.Model

	done chan struct{}
}

func NewInfo(styles config.Styles) *Info {
	return &Info{
		dialogBase: newDialogBase(styles),
		okButton:   button.New("Ok", true, styles),
		done:       make(chan struct{}, 1),
	}
}

func (m *Info) Init() tea.Cmd {
	return nil
}

func (m *Info) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.okButton, cmd = m.okButton.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && m.visible {
			m.visible = false
			m.done <- struct{}{}
		}
	}

	return m, cmd
}

func (m *Info) View() string {
	if !m.visible {
		return ""
	}

	text := m.styles.TextStyle.Render(m.text)

	okButton := m.okButton.View()

	output := lipgloss.JoinVertical(lipgloss.Center, text, okButton)

	return m.render(output)
}

func (m *Info) SetText(text string) {
	m.text = text
}

func (m *Info) Done() chan struct{} {
	return m.done
}
