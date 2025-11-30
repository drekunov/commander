package info

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/ui/config"
	"github.com/drekunov/gc/internal/ui/widgets/button"
)

type Model struct {
	title   string
	footer  string
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

func (m *Model) SetFooter(footer string) {
	m.footer = footer
}

func (m *Model) SetText(text string) {
	m.text = text
}

func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Model) SetWidth(width int) {
	m.width = width
}

func (m *Model) SetHeight(height int) {
	m.height = height
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

	header := m.headerView()

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

	footer := m.footerView()

	output = lipgloss.JoinVertical(lipgloss.Center, header, output, footer)

	return output
}

func (m *Model) headerView() string {
	borderStyle := config.Values.DialogBoxStyle.GetBorderStyle()
	header := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Center,
		" "+m.title+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top))

	header = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.TopLeft, header, borderStyle.TopRight)

	header = lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Render(header)

	return header
}

func (m *Model) footerView() string {
	borderStyle := config.Values.DialogBoxStyle.GetBorderStyle()
	footer := lipgloss.PlaceHorizontal(
		m.width,
		lipgloss.Center,
		" "+m.footer+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top))

	footer = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.BottomLeft, footer, borderStyle.BottomRight)

	footer = lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Render(footer)

	return footer
}
