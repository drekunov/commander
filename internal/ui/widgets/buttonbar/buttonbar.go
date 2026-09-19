// Package buttonbar provides a Norton-Commander-style bottom function-key
// button bar: a single full-width row of ten labeled buttons that are
// activated by function keys, mouse clicks, or keyboard navigation.
package buttonbar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
)

// Action identifies one of the ten function buttons on the bar.
type Action int

const (
	ActionHelp Action = iota + 1
	ActionMenu
	ActionView
	ActionEdit
	ActionCopy
	ActionRenMov
	ActionMkdir
	ActionDelete
	ActionPullDn
	ActionQuit
)

const (
	buttonCount = 10

	// naturalButtonWidth is a button's width when the terminal is wide enough:
	// a two-cell function-key number field plus a six-cell name field.
	naturalButtonWidth = 8
	numberFieldWidth   = 2
)

// actionDef describes a single bar button: its action identity and the label
// it is rendered under.
type actionDef struct {
	action Action
	label  string
	name   string
}

// actions lists the ten buttons left to right. The slice index maps to the
// matching function key: F1 activates actions[0], F2 actions[1], ..., F10
// actions[9].
var actions = []actionDef{
	{action: ActionHelp, label: "1Help", name: "Help"},
	{action: ActionMenu, label: "2Menu", name: "Menu"},
	{action: ActionView, label: "3View", name: "View"},
	{action: ActionEdit, label: "4Edit", name: "Edit"},
	{action: ActionCopy, label: "5Copy", name: "Copy"},
	{action: ActionRenMov, label: "6RenMov", name: "RenMov"},
	{action: ActionMkdir, label: "7Mkdir", name: "Mkdir"},
	{action: ActionDelete, label: "8Delete", name: "Delete"},
	{action: ActionPullDn, label: "9PullDn", name: "PullDn"},
	{action: ActionQuit, label: "10Quit", name: "Quit"},
}

// ActivateMsg reports that a bar button was activated. The consumer owns the
// action handlers and resolves the message (e.g. quit or run a mock).
type ActivateMsg struct {
	Action Action
}

// Model is the function-button bar widget.
type Model struct {
	styles config.Styles

	width int

	focused    bool
	focusedIdx int

	// pressed is the action activated by the most recent Update; it renders
	// in the pressed style for a single frame and is then cleared.
	pressed Action
}

// New creates an empty bar. Its width is set from window-size messages.
func New(styles config.Styles) *Model {
	return &Model{styles: styles}
}

// Label returns the button's display label, e.g. "5Copy", or "" when the
// action is not a known button.
func (a Action) Label() string {
	def, ok := actionDefFor(a)
	if !ok {
		return ""
	}

	return def.label
}

// Name returns the button's name without its function-key number, e.g. "Copy",
// or "" when the action is not a known button.
func (a Action) Name() string {
	def, ok := actionDefFor(a)
	if !ok {
		return ""
	}

	return def.name
}

// actionDefFor returns the definition for an action and reports whether the
// action is a known button.
func actionDefFor(a Action) (actionDef, bool) {
	idx := int(a) - 1
	if idx < 0 || idx >= len(actions) {
		return actionDef{}, false
	}

	return actions[idx], true
}

// Cmd returns a tea command that resolves to an ActivateMsg for the action.
func (a Action) Cmd() tea.Cmd {
	return func() tea.Msg {
		return ActivateMsg{Action: a}
	}
}

// ActionForKey returns the action bound to a function key.
func ActionForKey(key tea.KeyType) (Action, bool) {
	// Bubble Tea numbers F1..F20 with decreasing negative values, so the
	// offset from KeyF1 is measured as KeyF1-key.
	offset := int(tea.KeyF1 - key)
	if offset < 0 || offset >= buttonCount {
		return 0, false
	}

	return actions[offset].action, true
}

// Init implements tea.Model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	// Clear the one-frame pressed flash left by the previous activation.
	m.pressed = 0

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
	case tea.KeyMsg:
		return m, m.handleKey(msg)
	case tea.MouseMsg:
		return m, m.handleMouse(msg)
	}

	return m, nil
}

// View implements tea.Model.
func (m *Model) View() string {
	if m.width <= 0 {
		return ""
	}

	buttonWidth := m.buttonWidth()
	numberWidth := min(numberFieldWidth, buttonWidth)
	nameWidth := buttonWidth - numberWidth

	var row strings.Builder

	for idx, def := range actions {
		number := strings.TrimSuffix(def.label, def.name)
		row.WriteString(renderField(m.numberStyle(), number, numberWidth, true))
		row.WriteString(renderField(m.styleFor(idx, def.action), def.name, nameWidth, false))

		if idx < buttonCount-1 {
			if gap := m.gapWidth(idx); gap > 0 {
				row.WriteString(m.barBackground().Render(strings.Repeat(" ", gap)))
			}
		}
	}

	return row.String()
}

// renderField renders text padded to width cells with style, right-aligning it
// when alignRight is set and truncating it when it does not fit.
func renderField(style lipgloss.Style, text string, width int, alignRight bool) string {
	if width <= 0 {
		return ""
	}

	if lipgloss.Width(text) > width {
		text = ansi.Truncate(text, width, "")
	}

	pad := width - lipgloss.Width(text)

	if alignRight {
		return style.Render(strings.Repeat(" ", pad) + text)
	}

	return style.Render(text + strings.Repeat(" ", pad))
}

// SetWidth sets the bar's rendered width (the screen width).
func (m *Model) SetWidth(width int) {
	m.width = width
}

// Focused reports whether the bar holds keyboard focus.
func (m *Model) Focused() bool {
	return m.focused
}

// SetFocused marks whether the bar holds keyboard focus, which enables arrow
// navigation and Enter activation.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// ClearPressed drops the pressed flash so a button returns to its normal style
// once its activation has been dispatched, including while the dialog it opened
// is visible.
func (m *Model) ClearPressed() {
	m.pressed = 0
}

// handleKey routes key messages. Function keys always activate their button;
// arrows and Enter only act while the bar is focused.
func (m *Model) handleKey(msg tea.KeyMsg) tea.Cmd {
	if action, ok := ActionForKey(msg.Type); ok {
		return m.activate(action)
	}

	if !m.focused {
		return nil
	}

	switch msg.Type {
	case tea.KeyLeft:
		m.focusedIdx = (m.focusedIdx - 1 + buttonCount) % buttonCount
	case tea.KeyRight:
		m.focusedIdx = (m.focusedIdx + 1) % buttonCount
	case tea.KeyEnter:
		return m.activate(actions[m.focusedIdx].action)
	}

	return nil
}

// handleMouse activates the button under a left-button press on the bar row.
func (m *Model) handleMouse(msg tea.MouseMsg) tea.Cmd {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil
	}

	return m.activate(actions[m.actionIndexAtColumn(msg.X)].action)
}

// activate records the one-frame pressed state and returns an ActivateMsg cmd.
func (m *Model) activate(action Action) tea.Cmd {
	m.pressed = action

	return action.Cmd()
}

// barBackground returns the bar's background base: the injected menu-background
// style, falling back to the button-bar style for any attribute it leaves
// unset. It is the base for the padding between labels and for the number,
// name, and active/pressed styles, so the background applies in every state.
func (m *Model) barBackground() lipgloss.Style {
	return m.styles.MenuBackgroundStyle.Inherit(m.styles.ButtonBarStyle)
}

// numberStyle returns the style for a button's numeric prefix. It is drawn from
// the injected menu-number style on the bar background and falls back to that
// base for any attribute the menu-number style leaves unset. The number keeps
// this style on focused and pressed buttons.
func (m *Model) numberStyle() lipgloss.Style {
	return m.styles.MenuNumberStyle.Inherit(m.barBackground())
}

// styleFor picks the style for a button's name text. The pressed and focused
// buttons draw the pressed/active foreground and decoration on the bar strip's
// background; the remaining buttons draw the injected menu-label style on the
// strip, so their labels can be themed to match the cursor. The menu-label
// style falls back to the strip style for any attribute it leaves unset.
func (m *Model) styleFor(idx int, action Action) lipgloss.Style {
	switch {
	case action == m.pressed:
		return overlay(m.barBackground(), m.styles.PressedButtonStyle)
	case m.focused && idx == m.focusedIdx:
		return overlay(m.barBackground(), m.styles.ActiveButtonStyle)
	default:
		return m.styles.MenuLabelStyle.Inherit(m.barBackground())
	}
}

// overlay draws button's foreground and text decoration on top of bar's
// background, so labels sit on the themed strip instead of on their own fill.
func overlay(bar, button lipgloss.Style) lipgloss.Style {
	style := bar.Inherit(button)

	if fg := button.GetForeground(); fg != nil {
		style = style.Foreground(fg)
	}

	return style
}

// buttonWidth returns the width in cells of every button: the natural width
// when the terminal is wide enough, otherwise the widest equal width that fits.
func (m *Model) buttonWidth() int {
	width := m.width / buttonCount
	if width > naturalButtonWidth {
		return naturalButtonWidth
	}

	return width
}

// gapWidth returns the width in cells of the gap after the idx-th button. The
// columns left over after the ten buttons are split equally over the gaps
// between adjacent buttons; any remainder is spread over the leftmost gaps.
func (m *Model) gapWidth(idx int) int {
	gaps := buttonCount - 1
	if gaps <= 0 {
		return 0
	}

	leftover := m.width - buttonCount*m.buttonWidth()
	if leftover < 0 {
		return 0
	}

	base := leftover / gaps
	if idx < leftover%gaps {
		return base + 1
	}

	return base
}

// actionIndexAtColumn maps a screen column to the button whose span (its cells
// plus the following gap) contains the column.
func (m *Model) actionIndexAtColumn(col int) int {
	buttonWidth := m.buttonWidth()
	start := 0

	for idx := range actions {
		start += buttonWidth
		if idx < buttonCount-1 {
			start += m.gapWidth(idx)
		}

		if col < start {
			return idx
		}
	}

	return buttonCount - 1
}
