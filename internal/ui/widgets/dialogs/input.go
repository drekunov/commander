package dialogs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Input struct {
	dialogBase

	text string

	okButton *button.Model
	input    textinput.Model

	done chan struct{}
}

func NewInput(styles config.Styles, echo textinput.EchoMode) *Input {
	input := textinput.New()
	input.Placeholder = "Pikachu"
	input.Focus()
	input.CharLimit = 156
	input.Width = 20
	input.Prompt = ""
	input.EchoMode = echo

	return &Input{
		dialogBase: newDialogBase(styles),
		okButton:   button.New("Ok", true, styles),
		input:      input,
		done:       make(chan struct{}, 1),
	}
}

func (m *Input) Init() tea.Cmd {
	return nil
}

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds = make([]tea.Cmd, 0, 2)
		cmd  tea.Cmd
	)

	m.okButton, cmd = m.okButton.Update(msg)
	cmds = append(cmds, cmd)

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter && m.visible {
			m.visible = false
			m.done <- struct{}{}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Input) View() string {
	if !m.visible {
		return ""
	}

	text := m.styles.TextStyle.Render(m.text)

	input := m.styles.InputStyle.
		Render(m.input.View())

	okButton := m.okButton.View()

	output := lipgloss.JoinVertical(lipgloss.Center, text, input, okButton)

	return m.render(output)
}

func (m *Input) SetText(text string) {
	m.text = text
}

func (m *Input) Input() string {
	return m.input.Value()
}

func (m *Input) Done() chan struct{} {
	return m.done
}
