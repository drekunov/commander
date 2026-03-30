package wm

import (
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

const (
	minWindowWidth  = 10
	minWindowHeight = 5
)

// Manager is a window manager that holds multiple windows with z-index,
// move, and resize support. Each window wraps a tea.Model.
type Manager struct {
	windows       []*Window
	nextID        int
	width, height int
}

// New creates an empty window manager.
func New() *Manager {
	return &Manager{}
}

// SetSize sets the available screen area for the manager.
func (m *Manager) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// Width returns the manager's width.
func (m *Manager) Width() int { return m.width }

// Height returns the manager's height.
func (m *Manager) Height() int { return m.height }

// Add creates a new window containing the given tea.Model at position (x, y)
// with the given dimensions. Returns the window ID.
func (m *Manager) Add(content tea.Model, title string, x, y, w, h int) int {
	id := m.nextID
	m.nextID++

	win := NewWindow(id, content, x, y, w, h)
	win.Title = title
	win.ZIndex = m.topZIndex() + 1
	win.Focused = true

	// Unfocus all others.
	for _, ow := range m.windows {
		ow.Focused = false
	}

	m.windows = append(m.windows, win)
	return id
}

// Remove removes a window by ID.
func (m *Manager) Remove(id int) {
	for i, w := range m.windows {
		if w.ID == id {
			m.windows = append(m.windows[:i], m.windows[i+1:]...)
			break
		}
	}
	// Focus top-most remaining window.
	if top := m.topWindow(); top != nil {
		top.Focused = true
	}
}

// Focus brings a window to the front and focuses it.
func (m *Manager) Focus(id int) {
	for _, w := range m.windows {
		if w.ID == id {
			w.Focused = true
			w.ZIndex = m.topZIndex() + 1
		} else {
			w.Focused = false
		}
	}
}

// Move sets a window's position.
func (m *Manager) Move(id, x, y int) {
	if w := m.get(id); w != nil {
		w.X = x
		w.Y = y
	}
}

// Resize sets a window's dimensions, respecting minimums.
func (m *Manager) Resize(id, w, h int) {
	if win := m.get(id); win != nil {
		if w < minWindowWidth {
			w = minWindowWidth
		}
		if h < minWindowHeight {
			h = minWindowHeight
		}
		win.Width = w
		win.Height = h
	}
}

// SetZIndex sets a window's z-index.
func (m *Manager) SetZIndex(id, z int) {
	if w := m.get(id); w != nil {
		w.ZIndex = z
	}
}

// Get returns a window by ID, or nil if not found.
func (m *Manager) Get(id int) *Window {
	return m.get(id)
}

// Windows returns all windows sorted by z-index (back to front).
func (m *Manager) Windows() []*Window {
	sorted := make([]*Window, len(m.windows))
	copy(sorted, m.windows)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ZIndex < sorted[j].ZIndex
	})
	return sorted
}

// Init initializes all windows.
func (m *Manager) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, w := range m.windows {
		if cmd := w.Content.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// Update processes input messages. Mouse events handle drag/resize/focus;
// key events are routed only to the focused window.
func (m *Manager) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Other messages go to all windows.
	var cmds []tea.Cmd
	for _, w := range m.windows {
		var cmd tea.Cmd
		w.Content, cmd = w.Content.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

// View composites all visible windows onto a canvas, sorted by z-index.
func (m *Manager) View() string {
	if m.width == 0 || m.height == 0 {
		return ""
	}

	c := newCanvas(m.width, m.height)
	sorted := m.Windows()

	for _, w := range sorted {
		if !w.Visible {
			continue
		}
		rendered := m.renderWindow(w)
		c.stamp(w.X, w.Y, rendered)
	}

	return c.String()
}

// renderWindow draws a window frame with title bar and content.
func (m *Manager) renderWindow(w *Window) string {
	// Inner content area (subtract 2 for borders, 1 for title bar).
	innerW := w.Width - 2
	innerH := w.Height - 3
	if innerW < 1 {
		innerW = 1
	}
	if innerH < 1 {
		innerH = 1
	}

	// Render the tea.Model content.
	content := w.Content.View()

	// Build title bar.
	titleBar := m.buildTitleBar(w, innerW)

	// Build resize grip.
	grip := m.buildResizeGrip(innerW)

	// Frame the content with a border.
	border := config.Values.DialogBoxStyle.GetBorderStyle()
	borderColor := lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground())
	styledLeft := borderColor.Render(border.Left)
	styledRight := borderColor.Render(border.Right)

	framedContent := lipgloss.NewStyle().
		Width(innerW).
		Height(innerH).
		Render(content)

	// Compose: title bar + content + grip, with side borders.
	var b strings.Builder
	// Top border with title.
	b.WriteString(titleBar)
	b.WriteByte('\n')

	// Content lines with side borders.
	lines := strings.Split(framedContent, "\n")
	for _, line := range lines {
		b.WriteString(styledLeft)
		// Pad or clip line to innerW.
		lineW := lipgloss.Width(line)
		if lineW < innerW {
			line += strings.Repeat(" ", innerW-lineW)
		}
		b.WriteString(line)
		b.WriteString(styledRight)
		b.WriteByte('\n')
	}

	// Bottom border with resize grip.
	b.WriteString(grip)

	return b.String()
}

func (m *Manager) buildTitleBar(w *Window, innerW int) string {
	border := config.Values.DialogBoxStyle.GetBorderStyle()
	borderColor := config.Values.DialogBoxStyle.GetBorderBottomForeground()

	title := " " + w.Title + " "
	titleLen := lipgloss.Width(title)
	if titleLen > innerW {
		title = title[:innerW]
		titleLen = innerW
	}

	padding := innerW - titleLen
	left := padding / 2
	right := padding - left

	bar := border.TopLeft +
		strings.Repeat(border.Top, left) +
		title +
		strings.Repeat(border.Top, right) +
		border.TopRight

	style := lipgloss.NewStyle().Foreground(borderColor)
	if w.Focused {
		style = style.Bold(true)
	}

	return style.Render(bar)
}

func (m *Manager) buildResizeGrip(innerW int) string {
	border := config.Values.DialogBoxStyle.GetBorderStyle()

	gripRune := "◢"
	bottomLen := max(innerW-1, 0)

	grip := border.BottomLeft +
		strings.Repeat(border.Bottom, bottomLen) +
		gripRune +
		border.BottomRight

	return lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Render(grip)
}

func (m *Manager) handleMouse(msg tea.MouseMsg) tea.Cmd {
	x, y := msg.X, msg.Y

	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button != tea.MouseButtonLeft {
			return nil
		}
		// Check windows from front to back for click target.
		sorted := m.Windows()
		for i := len(sorted) - 1; i >= 0; i-- {
			w := sorted[i]
			if !w.Visible || !w.Contains(x, y) {
				continue
			}

			m.Focus(w.ID)

			if w.InTitleBar(x, y) {
				w.dragging = true
				w.dragOffX = x - w.X
				w.dragOffY = y - w.Y
				return nil
			}

			if w.InResizeGrip(x, y) {
				w.resizing = true
				w.resizeOrigW = w.Width
				w.resizeOrigH = w.Height
				w.resizeOrigX = x
				w.resizeOrigY = y
				return nil
			}

			// Click inside content area — translate coordinates and forward.
			localX := x - w.X - 1 // -1 for left border
			localY := y - w.Y - 1 // -1 for title bar
			localMsg := tea.MouseMsg{
				X:      localX,
				Y:      localY,
				Action: msg.Action,
				Button: msg.Button,
			}
			var cmd tea.Cmd
			w.Content, cmd = w.Content.Update(localMsg)
			return cmd
		}

	case tea.MouseActionMotion:
		for _, w := range m.windows {
			if w.dragging {
				w.X = x - w.dragOffX
				w.Y = y - w.dragOffY
				return nil
			}
			if w.resizing {
				dw := x - w.resizeOrigX
				dh := y - w.resizeOrigY
				newW := w.resizeOrigW + dw
				newH := w.resizeOrigH + dh
				m.Resize(w.ID, newW, newH)
				return nil
			}
		}

	case tea.MouseActionRelease:
		for _, w := range m.windows {
			w.dragging = false
			w.resizing = false
		}
	}

	return nil
}

func (m *Manager) handleKey(msg tea.KeyMsg) tea.Cmd {
	// Route key events only to the focused window.
	for _, w := range m.windows {
		if w.Focused && w.Visible {
			var cmd tea.Cmd
			w.Content, cmd = w.Content.Update(msg)
			return cmd
		}
	}
	return nil
}

// Helpers.

func (m *Manager) get(id int) *Window {
	for _, w := range m.windows {
		if w.ID == id {
			return w
		}
	}
	return nil
}

func (m *Manager) topZIndex() int {
	z := 0
	for _, w := range m.windows {
		if w.ZIndex > z {
			z = w.ZIndex
		}
	}
	return z
}

func (m *Manager) topWindow() *Window {
	var top *Window
	for _, w := range m.windows {
		if top == nil || w.ZIndex > top.ZIndex {
			top = w
		}
	}
	return top
}
