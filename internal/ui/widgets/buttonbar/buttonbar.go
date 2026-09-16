// Package buttonbar provides a Norton-Commander-style bottom function-key
// button bar: a single full-width row of ten labeled buttons that are
// activated by function keys, mouse clicks, or keyboard navigation.
package buttonbar

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

const buttonCount = 10

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

// Label returns the button's display label, e.g. "5Copy".
func (a Action) Label() string {
	return actions[int(a)-1].label
}

// Name returns the button's name without its function-key number, e.g. "Copy".
func (a Action) Name() string {
	return actions[int(a)-1].name
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

	var row strings.Builder

	for idx, def := range actions {
		segWidth := m.segmentWidth(idx)

		text := def.label
		if len(text) > segWidth {
			text = text[:segWidth]
		}

		text += strings.Repeat(" ", segWidth-len(text))

		style := colorize(m.styleFor(idx, def.action))
		row.WriteString(style.Render(text))
	}

	// The pressed style applies for the single frame that follows activation.
	m.pressed = 0

	return row.String()
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

// styleFor picks the button style for a segment: the pressed style for the
// button activated this frame, the active style for the focused button while
// the bar is focused, and the plain button style otherwise.
func (m *Model) styleFor(idx int, action Action) lipgloss.Style {
	switch {
	case action == m.pressed:
		return m.styles.PressedButtonStyle
	case m.focused && idx == m.focusedIdx:
		return m.styles.ActiveButtonStyle
	default:
		return m.styles.ButtonStyle
	}
}

// segmentWidth returns the width in cells of the idx-th button segment.
// The screen width is split into ten equal segments; any remainder is spread
// over the leftmost segments.
func (m *Model) segmentWidth(idx int) int {
	base := m.width / buttonCount
	if idx < m.width%buttonCount {
		return base + 1
	}

	return base
}

// actionIndexAtColumn maps a screen column to the button segment index.
func (m *Model) actionIndexAtColumn(col int) int {
	start := 0

	for i := range actions {
		start += m.segmentWidth(i)

		if col < start {
			return i
		}
	}

	return buttonCount - 1
}

// colorize returns a minimal style carrying the colors of the given style, so
// segments can be painted cell-for-cell at an exact width.
func colorize(style lipgloss.Style) lipgloss.Style {
	colored := lipgloss.NewStyle().
		Foreground(style.GetForeground()).
		Background(style.GetBackground())

	if style.GetBold() {
		colored = colored.Bold(true)
	}

	if style.GetUnderline() {
		colored = colored.Underline(true)
	}

	return colored
}
