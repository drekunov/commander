package panel

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
)

type Model struct {
	width   int
	height  int
	visible bool

	tableView table.Model
}

func NewPanel() *Model {
	return &Model{
		tableView: table.New(),
		visible:   true,
	}
}

// Focus activates the table cursor.
func (m *Model) Focus() {
	m.tableView.Focus()
}

// Blur hides the table cursor.
func (m *Model) Blur() {
	m.tableView.Blur()
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.tableView, cmd = m.tableView.Update(msg)

	return m, cmd
}

// View renders just the table content; the wm window provides the border frame.
func (m *Model) View() string {
	if !m.visible {
		return ""
	}

	m.tableView.SetWidth(m.width)
	m.tableView.SetHeight(m.height)

	return config.Values.TableStyle.Render(m.tableView.View())
}

func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Model) TableView() table.Model {
	return m.tableView
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
