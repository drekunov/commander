package panel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/evertras/bubble-table/table"
)

type Model struct {
	title   string
	footer  string
	width   int
	height  int
	visible bool

	tableView table.Model
}

func NewPanel() Model {
	tableView := table.New(nil)

	return Model{
		tableView: tableView,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	m.tableView, cmd = m.tableView.Update(msg)
	cmds = append(cmds, cmd)

	return *m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if !m.visible {
		return ""
	}

	tableView := config.Values.TextStyle.Render(m.tableView.View())

	output := config.Values.DialogBoxStyle.
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center, lipgloss.Center).
		BorderTop(false).
		BorderBottom(false).
		Render(tableView)

	header := m.headerView(lipgloss.Width(output) - 2)

	footer := m.footerView(lipgloss.Width(output) - 2)

	output = lipgloss.JoinVertical(lipgloss.Center, header, output, footer)

	return output
}

func (m *Model) SetTitle(title string) {
	m.title = title
}

func (m *Model) SetFooter(footer string) {
	m.footer = footer
}

func (m *Model) TableView() table.Model {
	return m.tableView
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

func (m *Model) SetTableView(tableView table.Model) {
	m.tableView = tableView
}

func (m *Model) headerView(width int) string {
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

func (m *Model) footerView(width int) string {
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
