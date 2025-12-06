package dialogs

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Input struct {
	title   string
	footer  string
	text    string
	width   int
	height  int
	visible bool

	okButton button.Model
	input    textinput.Model

	done chan struct{}
}

func NewInput() *Input {
	input := textinput.New()
	input.Placeholder = "Pikachu"
	input.Focus()
	input.CharLimit = 156
	input.Width = 20
	input.Prompt = ""

	return &Input{
		okButton: button.New("Ok", true),
		input:    input,
		done:     make(chan struct{}),
	}
}

func (m *Input) SetTitle(title string) {
	m.title = title
}

func (m *Input) SetFooter(footer string) {
	m.footer = footer
}

func (m *Input) SetText(text string) {
	m.text = text
}

func (m *Input) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Input) SetWidth(width int) {
	m.width = width
}

func (m *Input) SetHeight(height int) {
	m.height = height
}

func (m *Input) Input() string {
	return m.input.Value()
}

func (m *Input) Done() chan struct{} {
	return m.done
}

func (m *Input) Init() tea.Cmd {
	return nil
}

func (m *Input) Free() {
	close(m.done)
}

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.okButton, cmd = m.okButton.Update(msg)
	cmds = append(cmds, cmd)

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		if msg.Type == tea.KeyEnter {
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

	text := config.Values.TextStyle.Render(m.text)

	input := config.Values.InputStyle.
		Render(m.input.View())

	okButton := m.okButton.View()

	output := lipgloss.JoinVertical(lipgloss.Center, text, input, okButton)

	output = config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		BorderTop(false).
		BorderBottom(false).
		Render(output)

	header := m.headerView(lipgloss.Width(output) - 2)

	footer := m.footerView(lipgloss.Width(output) - 2)

	output = lipgloss.JoinVertical(lipgloss.Center, header, output, footer)

	output = lipgloss.Place(m.width, m.height, 0.1, 0.1, output)

	output = lipgloss.Place(m.width, m.height, 0.2, 0.2, output)

	return output
}

func (m *Input) headerView(width int) string {
	borderStyle := config.Values.DialogBoxStyle.GetBorderStyle()
	header := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		" "+m.title+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top))

	header = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.TopLeft, header, borderStyle.TopRight)

	header = lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Render(header)

	return header
}

func (m *Input) footerView(width int) string {
	borderStyle := config.Values.DialogBoxStyle.GetBorderStyle()
	footer := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		" "+m.footer+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top))

	footer = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.BottomLeft, footer, borderStyle.BottomRight)

	footer = lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Render(footer)

	return footer
}
