package mainform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

type Model struct {
	wm            *wm.Manager
	leftPanelID   int
	rightPanelID  int
	width, height int
}

func New() *Model {
	mgr := wm.New()

	left := panel.NewPanel()
	left.Focus()

	right := panel.NewPanel()

	leftID := mgr.Add(left, "", 0, 0, 40, 20)
	rightID := mgr.Add(right, "", 40, 0, 40, 20)

	// Restore focus to the left panel (Add focuses the last added window).
	mgr.Focus(leftID)

	return &Model{
		wm:           mgr,
		leftPanelID:  leftID,
		rightPanelID: rightID,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.wm.Init()
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyF10 {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layoutPanels()
	}

	cmd := m.wm.Update(msg)

	return m, cmd
}

func (m *Model) View() string {
	return m.wm.View()
}

// WM returns the underlying window manager (used for dialogs and overlays).
func (m *Model) WM() *wm.Manager {
	return m.wm
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

func (m *Model) Width() int  { return m.width }
func (m *Model) Height() int { return m.height }

// layoutPanels repositions and resizes the two panel windows to fill the screen.
func (m *Model) layoutPanels() {
	halfW := m.width / 2
	rightW := m.width - halfW

	m.wm.Move(m.leftPanelID, 0, 0)
	m.wm.Resize(m.leftPanelID, halfW, m.height)

	m.wm.Move(m.rightPanelID, halfW, 0)
	m.wm.Resize(m.rightPanelID, rightW, m.height)

	// Update panel inner dimensions (window border takes 2 cols, 3 rows).
	if p := m.panelFromWindow(m.leftPanelID); p != nil {
		p.SetWidth(halfW - 2)
		p.SetHeight(m.height - 3)
	}

	if p := m.panelFromWindow(m.rightPanelID); p != nil {
		p.SetWidth(rightW - 2)
		p.SetHeight(m.height - 3)
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
