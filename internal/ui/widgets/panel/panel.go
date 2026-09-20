package panel

import (
	"cmp"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/mattn/go-runewidth"
)

const (
	defaultColumnWidth = 12
	minFlexWidth       = 8
	parentLabel        = ".."
	parentIsDirValue   = "true"
	parentSizeValue    = "<DIR>"

	sortAscMarker  = "▲"
	sortDescMarker = "▼"

	attrName  = "Name"
	attrSize  = "Size"
	attrIsDir = "IsDir"
	attrKind  = "Kind"
)

// nameSGRSentinel is a private-use rune used to extract the opening ANSI
// sequence a style would emit for its text.
const nameSGRSentinel = "\uE000"

// columnSpec describes one visible table column declared by a header row.
type columnSpec struct {
	title string
	width int
	flex  bool
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

	// rowKinds holds the file kind of each display row, aligned with the order
	// produced by displayRowsFor (the synthetic parent row included), so the
	// render step can color each row's Name cell by kind.
	rowKinds []string

	// pendingSelect names the entry to place the cursor on when the next
	// listing arrives: set when ascending so the directory just left is
	// reselected, then cleared once the listing is applied.
	pendingSelect string

	// colSpecs holds the visible columns of the last delivered header row so
	// their widths can be recomputed when the panel is resized and a synthetic
	// parent (..) row can be built for any connector schema.
	colSpecs []columnSpec

	// sortCol and sortAsc hold the active sort column index and direction; the
	// default is the first column ascending.
	sortCol int
	sortAsc bool

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
		sortAsc:   true,
		styles:    styles,
	}

	model.tableView.SetStyles(model.tableStyles())
	model.tableView.KeyMap = navigationKeyMap()

	return model
}

func (m *Model) SetPanelID(id app.PanelID) {
	m.id = id
}

// ID returns the panel's identity.
func (m *Model) ID() app.PanelID {
	return m.id
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

// CursorVisible reports whether the panel currently renders its selection
// cursor, which is true only while the panel's table holds focus.
func (m *Model) CursorVisible() bool {
	return m.tableView.Focused()
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

	return m.colorizeKinds(m.styles.TableStyle.Render(m.tableView.View()))
}

// kindSGR returns the opening and closing ANSI sequences for a kind style. The
// closing sequence resets only the properties the style sets, so the table
// background and any cursor highlight survive around the colored name. It
// returns ok=false for the zero style or when colors are disabled.
func kindSGR(style lipgloss.Style) (string, string, bool) {
	rendered := style.Render(nameSGRSentinel)
	if rendered == nameSGRSentinel {
		return "", "", false
	}

	open := strings.TrimSuffix(rendered, nameSGRSentinel+"\x1b[0m")
	if open == "" || open == rendered {
		return "", "", false
	}

	var closeSeq strings.Builder

	if style.GetForeground() != (lipgloss.NoColor{}) {
		closeSeq.WriteString("\x1b[39m")
	}

	if style.GetBackground() != (lipgloss.NoColor{}) {
		closeSeq.WriteString("\x1b[49m")
	}

	if style.GetBold() {
		closeSeq.WriteString("\x1b[22m")
	}

	if style.GetUnderline() {
		closeSeq.WriteString("\x1b[24m")
	}

	return open, closeSeq.String(), true
}

// spliceCell wraps the cell at [left, left+width) of an ANSI-containing line
// with open and closeSeq, leaving every escape sequence already in the line in
// place so inherited styling continues through the cell.
func spliceCell(line string, left, width int, open, closeSeq string) string {
	begin := byteOffsetAt(line, left)

	end := byteOffsetAt(line, left+width)
	if begin >= end || end > len(line) {
		return line
	}

	return line[:begin] + open + line[begin:end] + closeSeq + line[end:]
}

// byteOffsetAt returns the byte index in text where the visible column col
// begins, counting ANSI escape sequences as zero width.
func byteOffsetAt(text string, col int) int {
	if col <= 0 {
		return 0
	}

	var (
		state  byte
		width  int
		pos    int
		parser = ansi.NewParser()
	)

	for pos < len(text) {
		_, cellWidth, read, newState := ansi.DecodeSequence(text[pos:], state, parser)
		if cellWidth > 0 && width+cellWidth > col {
			break
		}

		width += cellWidth
		state = newState
		pos += read
	}

	return pos
}

// plainRowLine reproduces the bubbles table's plain rendering of one row, so a
// rendered line can be matched back to the row that produced it.
func plainRowLine(row table.Row, cols []table.Column) string {
	var builder strings.Builder

	for idx, value := range row {
		if idx >= len(cols) || cols[idx].Width <= 0 {
			continue
		}

		width := cols[idx].Width
		text := runewidth.Truncate(value, width, "…")
		pad := width - lipgloss.Width(text)

		builder.WriteByte(' ')
		builder.WriteString(text)

		if pad > 0 {
			builder.WriteString(strings.Repeat(" ", pad))
		}

		builder.WriteByte(' ')
	}

	return builder.String()
}

// nameColumnLeft returns the visual column where the column at idx begins,
// including the per-cell left padding of every preceding column.
func nameColumnLeft(cols []table.Column, idx int) int {
	left := 1

	for j := range idx {
		left += cols[j].Width + 2
	}

	return left
}

func clampInt(value, low, high int) int {
	return min(max(value, low), high)
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
		m.columnsFor(data[0])
		m.rows = data[1:]
		m.sortRows()
	}

	// Drop the previous rows before the columns can change so SetColumns never
	// renders a stale row whose cell count no longer matches the new schema.
	m.tableView.SetRows(nil)
	m.applyColumnWidths()

	m.tableView.SetRows(m.displayRowsFor(m.rows))
	m.positionCursor(m.cursorForPendingSelect())
	m.pendingSelect = ""
}

func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

func (m *Model) SetWidth(width int) {
	m.width = width
	m.applyColumnWidths()
}

func (m *Model) SetHeight(height int) {
	m.height = height
}

// SortState returns the active sort column index and whether the sort is
// ascending.
func (m *Model) SortState() (int, bool) {
	return m.sortCol, m.sortAsc
}

// ColumnTitles returns the visible column titles in order, without the sort
// marker.
func (m *Model) ColumnTitles() []string {
	titles := make([]string, 0, len(m.colSpecs))
	for _, spec := range m.colSpecs {
		titles = append(titles, spec.title)
	}

	return titles
}

// SortBy sorts the panel by the named column. Choosing the current sort column
// inverts its direction; choosing another selects it ascending. The selected
// entry stays selected.
func (m *Model) SortBy(title string) {
	idx := m.columnIndex(title)
	if idx < 0 {
		return
	}

	if idx == m.sortCol {
		m.sortAsc = !m.sortAsc
	} else {
		m.sortCol = idx
		m.sortAsc = true
	}

	m.applySort()
}

// PreserveSelection marks the currently selected entry so the next listing
// delivered for the same directory restores it instead of falling back to the
// first row. A refresh uses it to keep the user's place.
func (m *Model) PreserveSelection() {
	m.pendingSelect = m.selectedEntryName()
}

// colorizeKinds paints the Name cell of each visible entry with the injected
// style for its file kind. It splices a foreground-only ANSI sequence into the
// already-rendered table, so the table's own width and truncation handling is
// untouched. Rows are matched by their plain rendered text rather than by line
// position, because the table does not expose its viewport scroll offset. The
// focused cursor row is left uncolored so the cursor style stays uniform.
func (m *Model) colorizeKinds(view string) string {
	cols := m.tableView.Columns()

	nameIdx, nameWidth := m.nameSpan(cols)
	if nameIdx < 0 || nameWidth <= 0 {
		return view
	}

	left := nameColumnLeft(cols, nameIdx)
	matches := m.rowLineMatches()

	lines := strings.Split(view, "\n")

	for idx := 1; idx < len(lines); idx++ {
		match, found := matches[strings.TrimRight(ansi.Strip(lines[idx]), " ")]
		if !found || match.conflict {
			continue
		}

		if match.isCursor && m.tableView.Focused() {
			continue
		}

		open, closeSeq, found := kindSGR(m.kindStyle(match.kind))
		if !found {
			continue
		}

		lines[idx] = spliceCell(lines[idx], left, nameWidth, open, closeSeq)
	}

	return strings.Join(lines, "\n")
}

// rowLineMatch describes the entry a rendered line belongs to.
type rowLineMatch struct {
	kind     string
	isCursor bool
	conflict bool
}

// rowLineMatches maps the plain rendered text of every row that can be visible
// to that row's kind and cursor state. The candidate window is the same one the
// bubbles table builds its content from, so every visible line is covered
// without relying on the table's internal scroll offset.
func (m *Model) rowLineMatches() map[string]rowLineMatch {
	rows := m.tableView.Rows()
	cols := m.tableView.Columns()
	cursor := m.tableView.Cursor()
	height := m.tableView.Height()

	start := clampInt(cursor-height, 0, cursor)
	end := clampInt(cursor+height, cursor, len(rows))

	matches := make(map[string]rowLineMatch, end-start)

	for idx := start; idx < end; idx++ {
		plain := strings.TrimRight(plainRowLine(rows[idx], cols), " ")
		info := rowLineMatch{kind: m.rowKindAt(idx), isCursor: idx == cursor}

		existing, found := matches[plain]
		if !found {
			matches[plain] = info

			continue
		}

		if existing.kind != info.kind || existing.isCursor != info.isCursor {
			existing.conflict = true
			matches[plain] = existing
		}
	}

	return matches
}

// rowKindAt returns the file kind of the display row at idx, or "" when the
// index is out of range.
func (m *Model) rowKindAt(idx int) string {
	if idx < 0 || idx >= len(m.rowKinds) {
		return ""
	}

	return m.rowKinds[idx]
}

// kindStyle returns the injected style for a file kind, or the zero style when
// the kind is absent or unknown.
func (m *Model) kindStyle(kind string) lipgloss.Style {
	switch kind {
	case app.FileKindDirectory:
		return m.styles.FileDirectoryStyle
	case app.FileKindExecutable:
		return m.styles.FileExecutableStyle
	case app.FileKindSymlink:
		return m.styles.FileSymlinkStyle
	case app.FileKindImage:
		return m.styles.FileImageStyle
	case app.FileKindArchive:
		return m.styles.FileArchiveStyle
	case app.FileKindSource:
		return m.styles.FileSourceStyle
	case app.FileKindConfig:
		return m.styles.FileConfigStyle
	default:
		return lipgloss.NewStyle()
	}
}

// nameSpan returns the index and width of the Name column, or -1 when it is not
// among the visible columns.
func (m *Model) nameSpan(cols []table.Column) (int, int) {
	idx := m.columnIndex(attrName)
	if idx < 0 || idx >= len(cols) {
		return -1, 0
	}

	return idx, cols[idx].Width
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

// columnsFor records the visible columns declared by a header row so their
// widths can be computed from the panel width. Hidden attributes are not
// columns; a header attribute's width, when positive, overrides the default
// column width.
func (m *Model) columnsFor(header app.AttributeList) {
	m.colSpecs = m.colSpecs[:0]

	for _, attr := range header {
		if attr.Hidden {
			continue
		}

		width := attr.Width
		if width <= 0 {
			width = defaultColumnWidth
		}

		m.colSpecs = append(m.colSpecs, columnSpec{
			title: attr.AttrName,
			width: width,
			flex:  attr.Flex,
		})
	}

	m.clampSortCol()
}

// clampSortCol resets the sort column to the first column when it no longer
// indexes a visible column. A listing may expose fewer columns than the one
// currently sorted, for example when the connector or directory schema changes.
func (m *Model) clampSortCol() {
	if m.sortCol < 0 || m.sortCol >= len(m.colSpecs) {
		m.sortCol = 0
		m.sortAsc = true
	}
}

// applyColumnWidths rebuilds the table columns from the recorded specs and the
// current panel width. The first flexible column takes the width left after the
// fixed columns and the per-cell padding, never below minFlexWidth; fixed
// columns keep their declared widths.
func (m *Model) applyColumnWidths() {
	if len(m.colSpecs) == 0 {
		return
	}

	padding := 2 * len(m.colSpecs)

	fixed := 0
	flexIdx := -1

	for idx, spec := range m.colSpecs {
		if spec.flex && flexIdx == -1 {
			flexIdx = idx

			continue
		}

		fixed += spec.width
	}

	cols := make([]table.Column, 0, len(m.colSpecs))

	for idx, spec := range m.colSpecs {
		width := spec.width
		if idx == flexIdx {
			width = max(m.width-padding-fixed, minFlexWidth)
		}

		cols = append(cols, table.Column{Title: m.headerTitle(idx), Width: width})
	}

	m.tableView.SetColumns(cols)
	m.contentWidth = contentWidthOf(cols)
}

// headerTitle returns the header title for the idx-th column, marking the
// active sort column with its direction.
func (m *Model) headerTitle(idx int) string {
	title := m.colSpecs[idx].title
	if idx != m.sortCol {
		return title
	}

	if m.sortAsc {
		return title + sortAscMarker
	}

	return title + sortDescMarker
}

// sortRows orders the rows by the active sort column and direction, keeping
// directories above files. The sort is stable, so entries that compare equal
// keep their previous order.
func (m *Model) sortRows() {
	if len(m.colSpecs) == 0 {
		return
	}

	sortCol := m.sortCol
	if sortCol < 0 || sortCol >= len(m.colSpecs) {
		sortCol = 0
	}

	column := m.colSpecs[sortCol].title

	sort.SliceStable(m.rows, func(left, right int) bool {
		leftDir := entryIsDir(m.rows[left])
		rightDir := entryIsDir(m.rows[right])

		if leftDir != rightDir {
			return leftDir
		}

		leftVal, _ := attrValue(m.rows[left], column)
		rightVal, _ := attrValue(m.rows[right], column)

		if m.sortAsc {
			return compareValues(leftVal, rightVal) < 0
		}

		return compareValues(leftVal, rightVal) > 0
	})
}

// applySort re-sorts the rows for the active mode and repositions the cursor on
// the selected entry.
func (m *Model) applySort() {
	m.pendingSelect = m.selectedEntryName()
	m.sortRows()
	m.tableView.SetRows(m.displayRowsFor(m.rows))
	m.positionCursor(m.cursorForPendingSelect())
	m.pendingSelect = ""
	m.applyColumnWidths()
}

// columnIndex returns the index of the column whose title matches name
// (case-insensitive, trimmed), or -1 when no column matches.
func (m *Model) columnIndex(name string) int {
	for idx, spec := range m.colSpecs {
		if strings.EqualFold(strings.TrimSpace(spec.title), strings.TrimSpace(name)) {
			return idx
		}
	}

	return -1
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

// compareValues orders two attribute values: two app.Size values compare by
// bytes, and every other value compares as case-insensitive text.
func compareValues(left, right any) int {
	leftSize, leftOK := left.(app.Size)
	rightSize, rightOK := right.(app.Size)

	if leftOK && rightOK {
		return cmp.Compare(int64(leftSize), int64(rightSize))
	}

	return strings.Compare(strings.ToLower(fmt.Sprint(left)), strings.ToLower(fmt.Sprint(right)))
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
		if attr.Hidden {
			continue
		}

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
// by one row per real entry, recording each row's file kind in the same order.
func (m *Model) displayRowsFor(entries []app.AttributeList) []table.Row {
	if len(m.tableView.Columns()) == 0 {
		m.setDefaultColumns()
	}

	m.rowKinds = m.rowKinds[:0]

	rows := make([]table.Row, 0, len(entries)+1)

	if m.showParent() {
		rows = append(rows, m.parentRow())
		m.rowKinds = append(m.rowKinds, app.FileKindDirectory)
	}

	for _, entry := range entries {
		rows = append(rows, rowFromAttrList(entry))
		m.rowKinds = append(m.rowKinds, entryKind(entry))
	}

	return rows
}

// entryKind returns an entry's hidden file kind, or "" when it carries none.
func entryKind(entry app.AttributeList) string {
	value, ok := attrValue(entry, attrKind)
	if !ok {
		return ""
	}

	return fmt.Sprint(value)
}

// setDefaultColumns installs columns for a listing delivered without a header
// (an empty directory), falling back to the columns seen on a previous listing
// or a single Name column.
func (m *Model) setDefaultColumns() {
	if len(m.colSpecs) == 0 {
		m.colSpecs = []columnSpec{{title: attrName, width: defaultColumnWidth}}
	}

	m.clampSortCol()
	m.applyColumnWidths()
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
	titles := m.ColumnTitles()
	if len(titles) == 0 {
		titles = []string{attrName}
	}

	row := make(table.Row, len(titles))

	nameIdx := max(titleIndex(titles, attrName), 0)
	row[nameIdx] = parentLabel

	if idx := titleIndex(titles, attrSize); idx >= 0 {
		row[idx] = parentSizeValue
	}

	if idx := titleIndex(titles, attrIsDir); idx >= 0 {
		row[idx] = parentIsDirValue
	}

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

// selectedEntryName returns the name of the entry under the cursor, or "" when
// the cursor is on the parent row.
func (m *Model) selectedEntryName() string {
	entry, ok := m.selectedEntry()
	if !ok {
		return ""
	}

	return entryName(entry)
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
