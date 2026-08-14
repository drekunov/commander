package mainform

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

type Model struct {
	wm            *wm.Manager
	leftPanelID   int
	rightPanelID  int
	width, height int
}

func New(styles config.Styles) *Model {
	mgr := wm.New(styles)

	left := panel.NewPanel(styles)
	left.Focus()

	right := panel.NewPanel(styles)

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
	case panel.DataMsg:
		m.setData(msg.Data)
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

// setData applies a listing to the left panel. Runs on the event loop.
func (m *Model) setData(data []app.AttributeList) {
	p := m.Panel()
	if p == nil {
		return
	}

	p.SetData(data)
}

// layoutPanels repositions and resizes the two panel windows to fill the screen.
func (m *Model) layoutPanels() {
	halfW := m.width / 2
	rightW := m.width - halfW

	m.wm.Move(m.leftPanelID, 0, 0)
	m.wm.Resize(m.leftPanelID, halfW, m.height)

	m.wm.Move(m.rightPanelID, halfW, 0)
	m.wm.Resize(m.rightPanelID, rightW, m.height)

	// Update panel inner dimensions (window frame takes the frame area).
	if p := m.panelFromWindow(m.leftPanelID); p != nil {
		p.SetWidth(halfW - wm.FrameCols)
		p.SetHeight(m.height - wm.FrameRows)
	}

	if p := m.panelFromWindow(m.rightPanelID); p != nil {
		p.SetWidth(rightW - wm.FrameCols)
		p.SetHeight(m.height - wm.FrameRows)
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
