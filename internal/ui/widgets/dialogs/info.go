package dialogs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Info struct {
	title   string
	footer  string
	text    string
	width   int
	height  int
	visible bool

	okButton button.Model

	done chan struct{}
}

func NewInfo() *Info {
	return &Info{
		okButton: button.New("Ok", true),
		done:     make(chan struct{}),
	}
}

func (m *Info) Init() tea.Cmd {
	return nil
}

func (m *Info) Free() {
	close(m.done)
}

func (m *Info) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds = make([]tea.Cmd, 0, 1)
		cmd  tea.Cmd
	)

	m.okButton, cmd = m.okButton.Update(msg)
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

func (m *Info) View() string {
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

func (m *Info) SetTitle(title string) {
	m.title = title
}

func (m *Info) SetFooter(footer string) {
	m.footer = footer
}

func (m *Info) SetText(text string) {
	m.text = text
}

func (m *Info) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Info) SetWidth(width int) {
	m.width = width
}

func (m *Info) SetHeight(height int) {
	m.height = height
}

func (m *Info) Done() chan struct{} {
	return m.done
}

func (m *Info) headerView(width int) string {
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

func (m *Info) footerView(width int) string {
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
