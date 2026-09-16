package panel

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
)

const (
	defaultColumnWidth = 12
	parentLabel        = ".."
	parentIsDirValue   = "true"
)

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

// DataMsg delivers a listing for one panel through the tea event loop so the
// table state is only ever touched from the Bubble Tea goroutine.
type DataMsg struct {
	Panel app.PanelID
	Path  string
	Data  []app.AttributeList
}

// NavigateMsg asks to change the current directory of a panel. It originates
// from the focused panel and is resolved by the ui layer through the app.
type NavigateMsg struct {
	Panel app.PanelID
	Dir   string
}

// Cmd returns a tea command that resolves to the navigation request.
func (n NavigateMsg) Cmd() tea.Cmd {
	return func() tea.Msg {
		return n
	}
}

type Model struct {
	id   app.PanelID
	dir  string
	rows []app.AttributeList

	// colTitles holds the column titles of the last delivered header row so a
	// synthetic parent (..) row can be built for any connector schema.
	colTitles []string

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
	model.tableView.KeyMap = navigationKeyMap()

	return model
}

func (m *Model) SetPanelID(id app.PanelID) {
	m.id = id
}

func (m *Model) Dir() string {
	return m.dir
}

func (m *Model) Focus() {
	m.tableView.Focus()
}

func (m *Model) Blur() {
	m.tableView.Blur()
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.tableView, cmd = m.tableView.Update(msg)

	if key, ok := msg.(tea.KeyMsg); ok {
		if nav := m.navigationCmd(key); nav != nil {
			return m, nav
		}
	}

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

// SetData rebuilds the table for dir from a listing whose first row is the
// column header and whose remaining rows are one directory entry each. When
// the directory has a parent, a synthetic ".." row leads the listing; the
// cursor is reset to the first row after a reload.
func (m *Model) SetData(dir string, data []app.AttributeList) {
	m.dir = dir

	if len(data) == 0 {
		m.rows = nil

		m.setDisplayRows()
		m.tableView.SetCursor(0)

		return
	}

	header := data[0]

	m.colTitles = make([]string, 0, len(header))
	cols := make([]table.Column, 0, len(header))

	for _, attr := range header {
		m.colTitles = append(m.colTitles, attr.AttrName)
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

	m.rows = data[1:]

	m.setDisplayRows()
	m.tableView.SetCursor(0)
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

// navigationKeyMap extends the default table keymap so Left jumps the cursor
// to the first row and Right steps one row down.
func navigationKeyMap() table.KeyMap {
	keyMap := table.DefaultKeyMap()

	keyMap.GotoTop.SetKeys(append([]string{"left"}, keyMap.GotoTop.Keys()...)...)
	keyMap.LineDown.SetKeys(append([]string{"right"}, keyMap.LineDown.Keys()...)...)

	return keyMap
}

func entryName(entry app.AttributeList) string {
	for _, attr := range entry {
		if attr.AttrName == "Name" {
			return fmt.Sprintf("%v", attr.AttrValue)
		}
	}

	return ""
}

func entryIsDir(entry app.AttributeList) bool {
	for _, attr := range entry {
		if attr.AttrName != "IsDir" {
			continue
		}

		switch value := attr.AttrValue.(type) {
		case bool:
			return value
		case string:
			return value == parentIsDirValue
		default:
			return false
		}
	}

	return false
}

func rowFromAttrList(attrList app.AttributeList) table.Row {
	rawData := make(table.Row, 0, len(attrList))

	for _, attr := range attrList {
		rawData = append(rawData, fmt.Sprintf("%v", attr.AttrValue))
	}

	return rawData
}

// showParent reports whether the current directory has a parent to ascend to,
// i.e. whether a synthetic ".." row leads the listing.
func (m *Model) showParent() bool {
	if m.dir == "" {
		return false
	}

	return filepath.Dir(m.dir) != m.dir
}

// setDisplayRows renders the synthetic parent row (when applicable) followed
// by one row per real entry.
func (m *Model) setDisplayRows() {
	if len(m.tableView.Columns()) == 0 {
		m.setDefaultColumns()
	}

	rows := make([]table.Row, 0, len(m.rows)+1)

	if m.showParent() {
		rows = append(rows, m.parentRow())
	}

	for _, entry := range m.rows {
		rows = append(rows, rowFromAttrList(entry))
	}

	m.tableView.SetRows(rows)
}

// setDefaultColumns installs columns for a listing delivered without a header
// (an empty directory), falling back to the titles seen on a previous listing
// or a single Name column.
func (m *Model) setDefaultColumns() {
	titles := m.colTitles
	if len(titles) == 0 {
		titles = []string{"Name"}
	}

	cols := make([]table.Column, 0, len(titles))
	for _, title := range titles {
		cols = append(cols, table.Column{
			Title: title,
			Width: defaultColumnWidth,
		})
	}

	m.colTitles = titles
	m.tableView.SetColumns(cols)
}

// parentRow builds the display row for the synthetic ".." entry from the
// delivered column titles.
func (m *Model) parentRow() table.Row {
	titles := m.colTitles
	if len(titles) == 0 {
		titles = []string{"Name"}
	}

	row := make(table.Row, len(titles))

	nameIdx := 0

	for idx, title := range titles {
		switch strings.ToLower(strings.TrimSpace(title)) {
		case "name":
			nameIdx = idx
		case "isdir":
			row[idx] = parentIsDirValue
		case "path":
			row[idx] = m.dir
		}
	}

	row[nameIdx] = parentLabel

	return row
}

// selectedEntryIndex maps the table cursor to the index of the selected real
// entry. It returns -1 when the cursor sits on the parent (..) row.
func (m *Model) selectedEntryIndex() int {
	idx := m.tableView.Cursor()
	if m.showParent() {
		if idx <= 0 {
			return -1
		}

		return idx - 1
	}

	return idx
}

// selectedEntry returns the attribute row under the table cursor.
func (m *Model) selectedEntry() (app.AttributeList, bool) {
	idx := m.selectedEntryIndex()
	if idx < 0 || idx >= len(m.rows) {
		return nil, false
	}

	return m.rows[idx], true
}

// navigationCmd turns Enter/Backspace on the focused panel into a directory
// change request. Enter descends into the selected directory (or ascends when
// the parent ".." row is selected); Backspace ascends to the parent.
func (m *Model) navigationCmd(msg tea.KeyMsg) tea.Cmd {
	switch msg.Type {
	case tea.KeyEnter:
		if m.selectedEntryIndex() == -1 {
			return m.ascendCmd()
		}

		entry, ok := m.selectedEntry()
		if !ok || !entryIsDir(entry) {
			return nil
		}

		return NavigateMsg{Panel: m.id, Dir: filepath.Join(m.dir, entryName(entry))}.Cmd()

	case tea.KeyBackspace:
		return m.ascendCmd()
	}

	return nil
}

// ascendCmd returns a navigation request to the parent directory, or nil when
// the panel is at the filesystem root (or has no directory yet).
func (m *Model) ascendCmd() tea.Cmd {
	if m.dir == "" {
		return nil
	}

	parent := filepath.Dir(m.dir)
	if parent == m.dir {
		return nil
	}

	return NavigateMsg{Panel: m.id, Dir: parent}.Cmd()
}
