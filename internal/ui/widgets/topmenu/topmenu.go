// Package topmenu provides a Norton-Commander-style top menu bar: a first-row
// line of pull-down menus that is activated with F9, navigated with the arrow
// keys or per-caption hotkey letters, and whose items dispatch the same actions
// as the function-button bar.
package topmenu

import (
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
)

// ActivateMsg reports a top-menu item activation. A non-zero Action means the
// item dispatches that function-button action; otherwise Label names a
// placeholder item that reports itself through the mock dialog.
type ActivateMsg struct {
	Label  string
	Action buttonbar.Action
}

// itemDef is one pull-down entry. An unset action marks a placeholder item.
type itemDef struct {
	label  string
	action buttonbar.Action
}

// menuDef is a top-level menu: its caption, its hotkey letter, and its items.
type menuDef struct {
	caption string
	hotkey  rune
	items   []itemDef
}

// labelFindFile is the Commands menu's first placeholder label.
const labelFindFile = "Find file"

// menus lists the top-level menus left to right. It is package data so adding a
// real action later is a table edit.
var menus = []menuDef{
	{
		caption: "Left",
		hotkey:  'L',
		items: []itemDef{
			{label: "Brief"},
			{label: "Full"},
			{label: "Info"},
			{label: "Tree"},
			{label: "Quick view"},
			{label: "Sort by..."},
			{label: "Filter..."},
		},
	},
	{
		caption: "Files",
		hotkey:  'F',
		items: []itemDef{
			{label: "View", action: buttonbar.ActionView},
			{label: "Edit", action: buttonbar.ActionEdit},
			{label: "Copy", action: buttonbar.ActionCopy},
			{label: "RenMov", action: buttonbar.ActionRenMov},
			{label: "Mkdir", action: buttonbar.ActionMkdir},
			{label: "Delete", action: buttonbar.ActionDelete},
		},
	},
	{
		caption: "Commands",
		hotkey:  'C',
		items: []itemDef{
			{label: labelFindFile},
			{label: "History"},
			{label: "Swap panels"},
			{label: "Compare directories"},
		},
	},
	{
		caption: "Options",
		hotkey:  'O',
		items: []itemDef{
			{label: "Configuration"},
			{label: "Save editor"},
			{label: "Save setup"},
		},
	},
	{
		caption: "Right",
		hotkey:  'R',
		items: []itemDef{
			{label: "Brief"},
			{label: "Full"},
			{label: "Info"},
			{label: "Tree"},
			{label: "Quick view"},
			{label: "Sort by..."},
			{label: "Filter..."},
		},
	},
}

// Model is the top menu bar widget.
type Model struct {
	styles config.Styles

	active  bool
	open    bool
	menuIdx int
	itemIdx int
}

// New creates an inactive top menu bar with the first menu highlighted on
// activation.
func New(styles config.Styles) *Model {
	return &Model{styles: styles}
}

// Activate makes the bar active, highlights the first menu, and leaves any
// pull-down closed.
func (m *Model) Activate() {
	m.active = true
	m.menuIdx = 0
	m.closeMenu()
}

// Deactivate drops the active state and closes any open pull-down.
func (m *Model) Deactivate() {
	m.active = false
	m.closeMenu()
}

// Toggle activates the bar when inactive and deactivates it otherwise.
func (m *Model) Toggle() {
	if m.active {
		m.Deactivate()

		return
	}

	m.Activate()
}

// Active reports whether the bar currently owns keyboard navigation.
func (m *Model) Active() bool {
	return m.active
}

// Open reports whether a pull-down is currently shown.
func (m *Model) Open() bool {
	return m.active && m.open
}

// Update routes a message while the bar handles input. F9 toggles the bar from
// any state; otherwise only an active bar consumes keys and every key it
// receives is considered handled so the focused panel does not react.
func (m *Model) Update(msg tea.Msg) (tea.Cmd, bool) {
	key, ok := msg.(tea.KeyMsg)
	if !ok {
		return nil, false
	}

	if key.Type == tea.KeyF9 {
		m.Toggle()

		return nil, true
	}

	if !m.active {
		return nil, false
	}

	return m.handleActiveKey(key), true
}

// View renders the full-width top row. Each caption keeps its hotkey letter in
// the hotkey style; the active caption's remaining text uses the active-caption
// style. The row is truncated to width and padded on the bar background.
func (m *Model) View(width int) string {
	if width <= 0 {
		return ""
	}

	var row strings.Builder

	for idx, menu := range menus {
		row.WriteString(m.captionSegment(idx, menu))
	}

	rendered := ansi.Truncate(row.String(), width, "")

	if pad := width - lipgloss.Width(rendered); pad > 0 {
		rendered += m.barStyle().Render(strings.Repeat(" ", pad))
	}

	return rendered
}

// PulldownHeight returns the number of screen rows the open pull-down covers
// below the top row, or 0 when no pull-down is open.
func (m *Model) PulldownHeight() int {
	if !m.Open() {
		return 0
	}

	return len(m.menuItems()) + 2
}

// PulldownBlock renders the open pull-down as a compact multi-line box and the
// column where its left edge starts. The block overlays the panels without
// blanking the columns around it, so no empty rows are produced. It returns
// "", 0 when no pull-down is open.
func (m *Model) PulldownBlock(width int) (string, int) {
	if width <= 0 || !m.Open() {
		return "", 0
	}

	items := m.menuItems()
	if len(items) == 0 {
		return "", 0
	}

	startX, boxWidth := m.pulldownBox(width)

	lines := make([]string, 0, len(items)+2)
	lines = append(lines, m.pulldownStyle().Render(borderLine(boxWidth, '┌', '─', '┐')))

	for idx, item := range items {
		lines = append(lines, m.pulldownItemRow(boxWidth, item.label, m.itemStyle(idx)))
	}

	lines = append(lines, m.pulldownStyle().Render(borderLine(boxWidth, '└', '─', '┘')))

	return strings.Join(lines, "\n"), startX
}

// HandleMouse routes a left-button press on the top row or an open pull-down.
// Clicking a caption activates the bar and opens that menu; clicking an item
// row activates the item. It reports whether the press was consumed and returns
// the item's activation command, if any. Presses outside the open box fall
// through so the panels keep receiving them.
func (m *Model) HandleMouse(msg tea.MouseMsg, width int) (tea.Cmd, bool) {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil, false
	}

	if msg.Y == 0 {
		m.clickCaption(msg.X)

		return nil, true
	}

	if m.PulldownHeight() == 0 || msg.Y > m.PulldownHeight() {
		return nil, false
	}

	idx := msg.Y - 2
	if idx < 0 || idx >= len(m.menuItems()) {
		return nil, false
	}

	startX, boxWidth := m.pulldownBox(width)
	if msg.X < startX || msg.X >= startX+boxWidth {
		return nil, false
	}

	m.itemIdx = idx

	return m.activateFocused(), true
}

// handleActiveKey applies one key to an active bar.
func (m *Model) handleActiveKey(key tea.KeyMsg) tea.Cmd {
	switch key.Type {
	case tea.KeyLeft:
		m.moveMenu(-1)
	case tea.KeyRight:
		m.moveMenu(1)
	case tea.KeyDown:
		m.moveDown()
	case tea.KeyUp:
		m.moveUp()
	case tea.KeyEnter:
		return m.handleEnter()
	case tea.KeyEsc:
		m.handleEscape()
	case tea.KeyRunes:
		m.selectByHotkey(key.Runes)
	}

	return nil
}

// handleEnter opens the active pull-down when closed and otherwise activates
// the focused item.
func (m *Model) handleEnter() tea.Cmd {
	if !m.open {
		m.openMenu()

		return nil
	}

	return m.activateFocused()
}

// handleEscape closes an open pull-down and otherwise deactivates the bar.
func (m *Model) handleEscape() {
	if m.open {
		m.closeMenu()

		return
	}

	m.Deactivate()
}

// menuItems returns the active menu's items, or nil when out of range.
func (m *Model) menuItems() []itemDef {
	if m.menuIdx < 0 || m.menuIdx >= len(menus) {
		return nil
	}

	return menus[m.menuIdx].items
}

// moveMenu moves the highlight by delta and keeps an open pull-down following
// the new menu.
func (m *Model) moveMenu(delta int) {
	count := len(menus)
	m.menuIdx = (m.menuIdx + delta + count) % count

	if m.open {
		m.itemIdx = 0
	}
}

// moveDown opens the active menu's pull-down when closed and otherwise moves
// the focus down within it.
func (m *Model) moveDown() {
	if !m.open {
		m.openMenu()

		return
	}

	if m.itemIdx < len(m.menuItems())-1 {
		m.itemIdx++
	}
}

// moveUp moves the focus up within an open pull-down; it does nothing when the
// pull-down is closed.
func (m *Model) moveUp() {
	if !m.open || m.itemIdx == 0 {
		return
	}

	m.itemIdx--
}

// selectByHotkey opens the first menu whose hotkey letter matches any of the
// typed runes, case-insensitively.
func (m *Model) selectByHotkey(runes []rune) {
	for _, typed := range runes {
		hotkey := unicode.ToUpper(typed)

		for idx, menu := range menus {
			if menu.hotkey == hotkey {
				m.menuIdx = idx
				m.openMenu()

				return
			}
		}
	}
}

// openMenu opens the active pull-down with its first item focused.
func (m *Model) openMenu() {
	m.open = true
	m.itemIdx = 0
}

// closeMenu closes the pull-down and resets its focus.
func (m *Model) closeMenu() {
	m.open = false
	m.itemIdx = 0
}

// activateFocused closes the pull-down, deactivates the bar, and returns a
// command carrying the focused item's activation.
func (m *Model) activateFocused() tea.Cmd {
	items := m.menuItems()
	if m.itemIdx < 0 || m.itemIdx >= len(items) {
		return nil
	}

	item := items[m.itemIdx]
	m.closeMenu()
	m.Deactivate()

	return func() tea.Msg {
		return ActivateMsg{Label: item.label, Action: item.action}
	}
}

// clickCaption activates the bar and opens the menu whose caption spans x.
func (m *Model) clickCaption(x int) {
	for idx, menu := range menus {
		start := m.captionStart(idx)
		if x >= start && x < start+m.captionWidth(menu) {
			m.active = true
			m.menuIdx = idx
			m.openMenu()

			return
		}
	}
}

// captionSegment renders one caption with its hotkey letter kept apart from the
// caption text. The caption's horizontal padding is rendered with the caption
// style, so the active highlight covers the padding as well as the text.
func (m *Model) captionSegment(idx int, menu menuDef) string {
	runes := []rune(menu.caption)
	if len(runes) == 0 {
		return ""
	}

	active := m.active && idx == m.menuIdx
	hotkey := m.hotkeyStyle(active).Padding(0, 0, 0, 0)
	caption := m.captionStyle(active).Padding(0, 0, 0, 0)

	left, right := m.captionPadding()

	return caption.Render(strings.Repeat(" ", left)) +
		hotkey.Render(string(runes[0])) +
		caption.Render(string(runes[1:])) +
		caption.Render(strings.Repeat(" ", right))
}

// captionPadding returns the horizontal padding shared by every caption. It is
// the larger of the active and inactive caption padding so the caption width
// stays constant across states; the padding comes from the theme's menu-caption
// classes.
func (m *Model) captionPadding() (int, int) {
	inactive := m.captionStyle(false)
	active := m.captionStyle(true)

	return max(inactive.GetPaddingLeft(), active.GetPaddingLeft()),
		max(inactive.GetPaddingRight(), active.GetPaddingRight())
}

// captionWidth is a caption's rendered width including its padding.
func (m *Model) captionWidth(menu menuDef) int {
	left, right := m.captionPadding()

	return left + len([]rune(menu.caption)) + right
}

// captionStart is the screen column where the idx-th caption begins.
func (m *Model) captionStart(idx int) int {
	start := 0

	for i := 0; i < idx && i < len(menus); i++ {
		start += m.captionWidth(menus[i])
	}

	return start
}

// pulldownInner is the width of the widest item label.
func (m *Model) pulldownInner() int {
	inner := 0

	for _, item := range m.menuItems() {
		if w := lipgloss.Width(item.label); w > inner {
			inner = w
		}
	}

	return inner
}

// pulldownBox returns the start column and width of the pull-down box, clamped
// to the available width.
func (m *Model) pulldownBox(width int) (int, int) {
	startX := m.captionStart(m.menuIdx)
	boxWidth := m.pulldownInner() + 4

	if startX > width-4 {
		startX = width - 4
	}

	if startX < 0 {
		startX = 0
	}

	if boxWidth > width-startX {
		boxWidth = width - startX
	}

	if boxWidth < 4 {
		boxWidth = 4
	}

	return startX, boxWidth
}

// pulldownItemRow renders one item row. The box borders use the pull-down frame
// style and only the item text (padded to the box's inner width) uses the item
// style, so the focused cursor never selects the frame.
func (m *Model) pulldownItemRow(boxWidth int, label string, itemStyle lipgloss.Style) string {
	frame := m.pulldownStyle()

	return frame.Render("│ ") +
		itemStyle.Render(m.pulldownItemText(boxWidth, label)) +
		frame.Render(" │")
}

// pulldownItemText truncates and pads a label to the box's inner width, without
// the surrounding frame.
func (m *Model) pulldownItemText(boxWidth int, label string) string {
	inner := boxWidth - 4

	text := ansi.Truncate(label, inner, "")

	pad := max(inner-lipgloss.Width(text), 0)

	return text + strings.Repeat(" ", pad)
}

// itemStyle returns the focused item's cursor style or the pull-down style for
// the other items.
func (m *Model) itemStyle(idx int) lipgloss.Style {
	if idx == m.itemIdx {
		return m.pulldownCursorStyle()
	}

	return m.pulldownStyle()
}

// borderLine builds a horizontal box border of the given width.
func borderLine(width int, left, mid, right rune) string {
	if width < 2 {
		width = 2
	}

	return string(left) + strings.Repeat(string(mid), width-2) + string(right)
}

// barStyle returns the top row background, falling back to the button bar.
func (m *Model) barStyle() lipgloss.Style {
	return m.styles.TopBarStyle.Inherit(m.styles.ButtonBarStyle)
}

// captionStyle returns the caption style for the active or inactive state.
func (m *Model) captionStyle(active bool) lipgloss.Style {
	if active {
		return m.styles.MenuCaptionActiveStyle.Inherit(m.barStyle())
	}

	return m.styles.MenuCaptionStyle.Inherit(m.barStyle())
}

// hotkeyStyle returns the hotkey-letter style. When the caption is active the
// letter keeps the hotkey foreground and decoration but sits on the active
// caption background, so the selection highlight spans the whole caption
// instead of leaving the first letter unselected.
func (m *Model) hotkeyStyle(active bool) lipgloss.Style {
	if active {
		style := m.captionStyle(true)
		hotkey := m.styles.MenuHotkeyStyle

		if fg := hotkey.GetForeground(); fg != nil {
			style = style.Foreground(fg)
		}

		if hotkey.GetBold() {
			style = style.Bold(true)
		}

		return style
	}

	return m.styles.MenuHotkeyStyle.Inherit(m.captionStyle(false))
}

// pulldownStyle returns the pull-down background, falling back to the menu
// background and then the bar.
func (m *Model) pulldownStyle() lipgloss.Style {
	return m.styles.PulldownStyle.Inherit(m.styles.MenuBackgroundStyle).Inherit(m.barStyle())
}

// pulldownCursorStyle returns the focused-item style, falling back to the
// pull-down background.
func (m *Model) pulldownCursorStyle() lipgloss.Style {
	return m.styles.PulldownCursorStyle.Inherit(m.pulldownStyle())
}
