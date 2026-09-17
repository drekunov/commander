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

	attrName  = "Name"
	attrIsDir = "IsDir"
)

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

	// pendingSelect names the entry to place the cursor on when the next
	// listing arrives: set when ascending so the directory just left is
	// reselected, then cleared once the listing is applied.
	pendingSelect string

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

	model.tableView.SetStyles(model.tableStyles())
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
	m.tableView.SetStyles(m.tableStyles())

	return m.styles.TableStyle.Render(m.tableView.View())
}

// SetData rebuilds the table for dir from a listing whose first row is the
// column header and whose remaining rows are one directory entry each. When
// the directory has a parent, a synthetic ".." row leads the listing. The
// cursor is restored to the directory that was just left when ascending, and
// placed on the first row otherwise. The listing is scrolled so the selected
// row is visible, including the restored directory after an ascent.
func (m *Model) SetData(dir string, data []app.AttributeList) {
	m.dir = dir

	if len(data) == 0 {
		m.rows = nil
	} else {
		cols := m.columnsFor(data[0])
		m.tableView.SetColumns(cols)
		m.contentWidth = contentWidthOf(cols)
		m.rows = data[1:]
	}

	m.tableView.SetRows(m.displayRowsFor(m.rows))
	m.positionCursor(m.cursorForPendingSelect())
	m.pendingSelect = ""
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

// positionCursor places the table cursor on display row idx and scrolls the
// listing so that row is visible. GotoTop resets both the cursor and the
// viewport offset; MoveDown then advances the cursor while keeping it in view.
// This keeps the restored directory visible after an ascent even when the
// parent listing is taller than the panel.
func (m *Model) positionCursor(idx int) {
	m.tableView.GotoTop()
	m.tableView.MoveDown(idx)
}

// tableStyles returns the table styles for the current focus state. The cursor
// is drawn from the injected cursor style, and the header row is given the
// injected table background so it stays continuous with the listing in both
// focus states. The bubbles table highlights the cursor row regardless of
// focus, so a blurred panel renders that row unstyled to hide its cursor. A
// focused panel spans the cursor bar across the full panel width so the
// viewport's horizontal padding stays painted with the cursor background
// instead of leaking as unstyled cells after the selected row's style reset.
func (m *Model) tableStyles() table.Styles {
	styles := table.DefaultStyles()
	styles.Header = styles.Header.Background(m.styles.TableStyle.GetBackground())
	styles.Selected = m.styles.CursorStyle

	if !m.tableView.Focused() {
		styles.Selected = lipgloss.NewStyle()

		return styles
	}

	if padRight := m.width - m.contentWidth; padRight > 0 {
		styles.Selected = styles.Selected.PaddingRight(padRight)
	}

	return styles
}

// columnsFor builds the table columns from a header row and records their
// titles so a synthetic parent row can be built for any connector schema.
func (m *Model) columnsFor(header app.AttributeList) []table.Column {
	m.colTitles = make([]string, 0, len(header))
	cols := make([]table.Column, 0, len(header))

	for _, attr := range header {
		m.colTitles = append(m.colTitles, attr.AttrName)
		cols = append(cols, table.Column{
			Title: attr.AttrName,
			Width: defaultColumnWidth,
		})
	}

	return cols
}

// contentWidthOf returns the natural rendered width of a set of columns,
// including one cell of padding on each side of every cell.
func contentWidthOf(cols []table.Column) int {
	width := 0
	for _, col := range cols {
		width += col.Width
	}

	return width + 2*len(cols)
}

// navigationKeyMap extends the default table keymap so Left jumps the cursor
// to the first row and Right jumps the cursor to the last row.
func navigationKeyMap() table.KeyMap {
	keyMap := table.DefaultKeyMap()

	keyMap.GotoTop.SetKeys(append([]string{"left"}, keyMap.GotoTop.Keys()...)...)
	keyMap.GotoBottom.SetKeys(append([]string{"right"}, keyMap.GotoBottom.Keys()...)...)

	return keyMap
}

// attrValue returns the value of the named attribute using a case-insensitive,
// whitespace-trimmed name match so every connector schema is read alike.
func attrValue(entry app.AttributeList, name string) (any, bool) {
	for _, attr := range entry {
		if strings.EqualFold(strings.TrimSpace(attr.AttrName), name) {
			return attr.AttrValue, true
		}
	}

	return nil, false
}

func entryName(entry app.AttributeList) string {
	value, ok := attrValue(entry, attrName)
	if !ok {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

func entryIsDir(entry app.AttributeList) bool {
	value, ok := attrValue(entry, attrIsDir)
	if !ok {
		return false
	}

	switch value := value.(type) {
	case bool:
		return value
	case string:
		return strings.EqualFold(strings.TrimSpace(value), parentIsDirValue)
	default:
		return false
	}
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

// cursorForPendingSelect returns the display-row index of the entry named by
// pendingSelect, offset by the synthetic parent row when it is shown. It
// returns 0 when there is nothing to reselect or the entry is absent.
func (m *Model) cursorForPendingSelect() int {
	if m.pendingSelect == "" {
		return 0
	}

	offset := 0
	if m.showParent() {
		offset = 1
	}

	for idx, entry := range m.rows {
		if entryName(entry) == m.pendingSelect {
			return idx + offset
		}
	}

	return 0
}

// displayRowsFor renders the synthetic parent row (when applicable) followed
// by one row per real entry.
func (m *Model) displayRowsFor(entries []app.AttributeList) []table.Row {
	if len(m.tableView.Columns()) == 0 {
		m.setDefaultColumns()
	}

	rows := make([]table.Row, 0, len(entries)+1)

	if m.showParent() {
		rows = append(rows, m.parentRow())
	}

	for _, entry := range entries {
		rows = append(rows, rowFromAttrList(entry))
	}

	return rows
}

// setDefaultColumns installs columns for a listing delivered without a header
// (an empty directory), falling back to the titles seen on a previous listing
// or a single Name column.
func (m *Model) setDefaultColumns() {
	titles := m.colTitles
	if len(titles) == 0 {
		titles = []string{attrName}
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

// titleIndex returns the index of the column whose title matches name
// (case-insensitive, trimmed), or -1 when no column matches.
func titleIndex(titles []string, name string) int {
	for idx, title := range titles {
		if strings.EqualFold(strings.TrimSpace(title), name) {
			return idx
		}
	}

	return -1
}

// parentRow builds the display row for the synthetic ".." entry from the
// delivered column titles.
func (m *Model) parentRow() table.Row {
	titles := m.colTitles
	if len(titles) == 0 {
		titles = []string{attrName}
	}

	row := make(table.Row, len(titles))

	nameIdx := max(titleIndex(titles, attrName), 0)

	if idx := titleIndex(titles, attrIsDir); idx >= 0 {
		row[idx] = parentIsDirValue
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

		m.pendingSelect = ""

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

	m.pendingSelect = filepath.Base(m.dir)

	return NavigateMsg{Panel: m.id, Dir: parent}.Cmd()
}
