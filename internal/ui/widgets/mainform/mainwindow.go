package mainform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

const functionBarHeight = 1

type Model struct {
	wm            *wm.Manager
	leftPanelID   int
	rightPanelID  int
	width, height int

	bar *buttonbar.Model
}

func New(styles config.Styles) *Model {
	mgr := wm.New(styles)

	left := panel.NewPanel(styles)
	left.SetPanelID(app.PanelLeft)

	right := panel.NewPanel(styles)
	right.SetPanelID(app.PanelRight)

	leftID := mgr.Add(left, "", 0, 0, 40, 20)
	rightID := mgr.Add(right, "", 40, 0, 40, 20)

	// Restore focus to the left panel (Add focuses the last added window).
	mgr.Focus(leftID)

	model := &Model{
		wm:           mgr,
		leftPanelID:  leftID,
		rightPanelID: rightID,
		bar:          buttonbar.New(styles),
	}
	model.syncPanelFocus()

	return model
}

func (m *Model) Init() tea.Cmd {
	return m.wm.Init()
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.applyWindowSize()

		return m, nil

	case tea.KeyMsg:
		if cmd, handled := m.handleKey(msg); handled {
			return m, cmd
		}

	case tea.MouseMsg:
		if cmd, handled := m.handleMouse(msg); handled {
			return m, cmd
		}

	case panel.DataMsg:
		m.applyData(msg)

		return m, nil
	}

	cmd := m.wm.Update(msg)
	m.syncPanelFocus()

	return m, cmd
}

func (m *Model) View() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}

	bar := m.bar.View()

	if m.height <= functionBarHeight {
		return bar
	}

	canvas := m.wm.View()
	if canvas == "" {
		return bar
	}

	return canvas + "\n" + bar
}

// WM returns the underlying window manager (used for dialogs and overlays).
func (m *Model) WM() *wm.Manager {
	return m.wm
}

// Bar returns the function-button bar widget.
func (m *Model) Bar() *buttonbar.Model {
	return m.bar
}

// Panel returns the left (first) panel.
func (m *Model) Panel() *panel.Model {
	win := m.wm.Get(m.leftPanelID)
	if win == nil {
		return nil
	}

	p, _ := win.Content.(*panel.Model)

	return p
}

// FocusedPanel returns the panel window that currently holds focus, or nil when
// no panel is focused.
func (m *Model) FocusedPanel() *panel.Model {
	for _, win := range m.wm.Windows() {
		if !win.Focused {
			continue
		}

		if p, ok := win.Content.(*panel.Model); ok {
			return p
		}
	}

	return nil
}

func (m *Model) Width() int  { return m.width }
func (m *Model) Height() int { return m.height }

// dialogOpen reports whether a dialog window (not a panel) currently holds
// focus.
func (m *Model) dialogOpen() bool {
	id := m.wm.FocusedWindowID()

	return id >= 0 && id != m.leftPanelID && id != m.rightPanelID
}

// handleKey routes keys between the function-button bar and the window
// manager. Function keys activate their button before windows see them;
// arrows and Enter drive the bar only while it holds keyboard focus. While a
// modal dialog is open, F1-F9 are ignored and other keys go to the dialog.
func (m *Model) handleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if msg.Type == tea.KeyF10 {
		return tea.Quit, true
	}

	if m.dialogOpen() {
		if _, ok := buttonbar.ActionForKey(msg.Type); ok {
			return nil, true
		}

		return nil, false
	}

	if _, ok := buttonbar.ActionForKey(msg.Type); ok {
		_, cmd := m.bar.Update(msg)

		return cmd, true
	}

	if m.bar.Focused() {
		switch msg.Type {
		case tea.KeyLeft, tea.KeyRight, tea.KeyEnter:
			_, cmd := m.bar.Update(msg)

			return cmd, true

		case tea.KeyTab:
			m.bar.SetFocused(false)
			m.wm.CycleFocus()
			m.syncPanelFocus()

			return nil, true
		}
	}

	return nil, false
}

// handleMouse intercepts presses on the bottom bar row, focusing the bar and
// activating the clicked button. Presses anywhere else clear bar focus so the
// window manager keeps handling them.
func (m *Model) handleMouse(msg tea.MouseMsg) (tea.Cmd, bool) {
	if msg.Action != tea.MouseActionPress || msg.Button != tea.MouseButtonLeft {
		return nil, false
	}

	if m.dialogOpen() {
		// A modal dialog owns the pointer: clicks outside it are ignored.
		return nil, true
	}

	if msg.Y == m.height-functionBarHeight {
		m.bar.SetFocused(true)
		m.bar.SetWidth(m.width)

		_, cmd := m.bar.Update(msg)

		return cmd, true
	}

	m.bar.SetFocused(false)

	return nil, false
}

// applyWindowSize sizes the bar and the window-manager area below it.
func (m *Model) applyWindowSize() {
	m.bar.SetWidth(m.width)
	m.wm.SetSize(m.width, m.wmHeight())
	m.layoutPanels()
}

// wmHeight returns the canvas height left for windows above the bar row.
func (m *Model) wmHeight() int {
	height := m.height - functionBarHeight
	if height < 1 {
		return 0
	}

	return height
}

// applyData routes a delivered listing to the addressed panel window and
// reflects the new directory in its title.
func (m *Model) applyData(msg panel.DataMsg) {
	var winID int

	switch msg.Panel {
	case app.PanelLeft:
		winID = m.leftPanelID
	case app.PanelRight:
		winID = m.rightPanelID
	default:
		return
	}

	win := m.wm.Get(winID)
	if win == nil {
		return
	}

	panelModel, ok := win.Content.(*panel.Model)
	if !ok {
		return
	}

	win.Title = msg.Path
	panelModel.SetData(msg.Path, msg.Data)
}

// layoutPanels repositions and resizes the two panel windows to fill the
// area above the function-button bar.
func (m *Model) layoutPanels() {
	areaHeight := m.wmHeight()
	if areaHeight < 1 {
		return
	}

	halfW := m.width / 2
	rightW := m.width - halfW

	m.wm.Move(m.leftPanelID, 0, 0)
	m.wm.Resize(m.leftPanelID, halfW, areaHeight)

	m.wm.Move(m.rightPanelID, halfW, 0)
	m.wm.Resize(m.rightPanelID, rightW, areaHeight)

	// Update panel inner dimensions (window frame takes the frame area).
	if p := m.panelFromWindow(m.leftPanelID); p != nil {
		p.SetWidth(halfW - wm.FrameCols)
		p.SetHeight(areaHeight - wm.FrameRows)
	}

	if p := m.panelFromWindow(m.rightPanelID); p != nil {
		p.SetWidth(rightW - wm.FrameCols)
		p.SetHeight(areaHeight - wm.FrameRows)
	}
}

func (m *Model) panelFromWindow(winID int) *panel.Model {
	win := m.wm.Get(winID)
	if win == nil {
		return nil
	}

	p, _ := win.Content.(*panel.Model)

	return p
}

// syncPanelFocus mirrors each panel window's focus state onto its table
// cursor so only the focused panel renders a selection cursor.
func (m *Model) syncPanelFocus() {
	for _, id := range []int{m.leftPanelID, m.rightPanelID} {
		win := m.wm.Get(id)
		if win == nil {
			continue
		}

		panelModel, ok := win.Content.(*panel.Model)
		if !ok {
			continue
		}

		if win.Focused {
			panelModel.Focus()
		} else {
			panelModel.Blur()
		}
	}
}
