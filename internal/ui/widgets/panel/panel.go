package panel

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
)

const defaultColumnWidth = 12

// ncTableStyles returns the table styles matching the Norton Commander look:
// the cursor is a black-on-gray bar instead of bubbles' default magenta.
func ncTableStyles() table.Styles {
	styles := table.DefaultStyles()

	styles.Selected = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#C0C0C0"))

	return styles
}

// DataMsg delivers a listing to the panel through the tea event loop so the
// table state is only ever touched from the Bubble Tea goroutine.
type DataMsg struct {
	Data []app.AttributeList
}

type Model struct {
	width   int
	height  int
	visible bool

	styles config.Styles

	// contentWidth is the natural rendered width of the table rows (column
	// widths plus the per-cell horizontal padding).
	contentWidth int

	tableView table.Model
}

func NewPanel(styles config.Styles) *Model {
	model := &Model{
		tableView: table.New(),
		visible:   true,
		styles:    styles,
	}

	model.tableView.SetStyles(ncTableStyles())

	return model
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

	// Span the cursor bar across the full panel width so the table viewport's
	// horizontal padding stays gray instead of leaking as unstyled cells after
	// the selected row's style reset.
	styles := ncTableStyles()
	if padRight := m.width - m.contentWidth; padRight > 0 {
		styles.Selected = styles.Selected.PaddingRight(padRight)
	}

	m.tableView.SetStyles(styles)

	return m.styles.TableStyle.Render(m.tableView.View())
}

// SetData rebuilds the table from a listing whose first row is the column
// header and whose remaining rows are one directory entry each.
func (m *Model) SetData(data []app.AttributeList) {
	if len(data) == 0 {
		return
	}

	header := data[0]

	cols := make([]table.Column, 0, len(header))
	for _, attr := range header {
		cols = append(cols, table.Column{
			Title: attr.AttrName,
			Width: defaultColumnWidth,
		})
	}

	m.tableView.SetColumns(cols)

	contentWidth := 0
	for _, col := range cols {
		contentWidth += col.Width
	}

	// Each cell carries one column of padding on the left and right.
	m.contentWidth = contentWidth + 2*len(cols)

	rows := make([]table.Row, 0, len(data)-1)
	for _, row := range data[1:] {
		rows = append(rows, rowFromAttrList(row))
	}

	m.tableView.SetRows(rows)
}

func rowFromAttrList(attrList app.AttributeList) table.Row {
	rawData := make(table.Row, 0, len(attrList))

	for _, attr := range attrList {
		rawData = append(rawData, fmt.Sprintf("%v", attr.AttrValue))
	}

	return rawData
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
