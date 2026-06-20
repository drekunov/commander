package wm

import (
	"slices"
	"sort"
	"strings"
	"sync"

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
	mu            sync.Mutex
	windows       []*Window
	nextID        int
	width, height int
}

// New creates an empty window manager.
func New() *Manager {
	return &Manager{}
}

// SetSize sets the available screen area for the manager.
func (m *Manager) SetSize(width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.width = width
	m.height = height
}

// Width returns the manager's width.
func (m *Manager) Width() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.width
}

// Height returns the manager's height.
func (m *Manager) Height() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.height
}

// Add creates a new window containing the given tea.Model at position (posX, posY)
// with the given dimensions. Returns the window ID.
func (m *Manager) Add(content tea.Model, title string, posX, posY, width, height int) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	winID := m.nextID
	m.nextID++

	win := NewWindow(winID, content, posX, posY, width, height)
	win.Title = title
	win.ZIndex = m.topZIndex() + 1
	win.Focused = true

	// Unfocus all others.
	for _, other := range m.windows {
		other.Focused = false
	}

	m.windows = append(m.windows, win)

	return winID
}

// Remove removes a window by ID.
func (m *Manager) Remove(winID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, win := range m.windows {
		if win.ID == winID {
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
func (m *Manager) Focus(winID int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.focus(winID)
}

// Move sets a window's position.
func (m *Manager) Move(winID, posX, posY int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if win := m.get(winID); win != nil {
		win.X = posX
		win.Y = posY
	}
}

// Resize sets a window's dimensions, respecting minimums.
func (m *Manager) Resize(winID, width, height int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.resize(winID, width, height)
}

// SetZIndex sets a window's z-index.
func (m *Manager) SetZIndex(winID, zIndex int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if win := m.get(winID); win != nil {
		win.ZIndex = zIndex
	}
}

// Get returns a window by ID, or nil if not found.
func (m *Manager) Get(winID int) *Window {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.get(winID)
}

// Windows returns all windows sorted by z-index (back to front).
func (m *Manager) Windows() []*Window {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.sortedWindows()
}

// Init initializes all windows.
func (m *Manager) Init() tea.Cmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	var cmds []tea.Cmd

	for _, win := range m.windows {
		if cmd := win.Content.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// Update processes input messages. Mouse events handle drag/resize/focus;
// Tab cycles window focus; key events are routed only to the focused window.
func (m *Manager) Update(msg tea.Msg) tea.Cmd {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		return nil
	case tea.MouseMsg:
		return m.handleMouse(msg)
	case tea.KeyMsg:
		if msg.Type == tea.KeyTab {
			m.cycleFocus()

			return nil
		}

		return m.handleKey(msg)
	}

	// Other messages go to all windows.
	var cmds []tea.Cmd

	for _, win := range m.windows {
		var cmd tea.Cmd

		win.Content, cmd = win.Content.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return tea.Batch(cmds...)
}

// CycleFocus moves focus to the next window in z-order.
func (m *Manager) CycleFocus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cycleFocus()
}

// View composites all visible windows onto a canvas, sorted by z-index.
// Returns "" when there is nothing to overlay so callers can skip compositing.
func (m *Manager) View() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.width == 0 || m.height == 0 {
		return ""
	}

	sorted := m.sortedWindows()

	// Collect only visible windows; skip compositing entirely if there are none.
	var visible []*Window

	for _, win := range sorted {
		if win.Visible {
			visible = append(visible, win)
		}
	}

	if len(visible) == 0 {
		return ""
	}

	cvs := newCanvas(m.width, m.height)

	for _, win := range visible {
		rendered := m.renderWindow(win)
		cvs.stamp(win.X, win.Y, rendered)
	}

	return cvs.String()
}

// focus is the lock-free internal version of Focus.
func (m *Manager) focus(winID int) {
	for _, win := range m.windows {
		if win.ID == winID {
			win.Focused = true
			win.ZIndex = m.topZIndex() + 1
		} else {
			win.Focused = false
		}
	}
}

// resize is the lock-free internal version of Resize.
func (m *Manager) resize(winID, width, height int) {
	win := m.get(winID)
	if win == nil {
		return
	}

	if width < minWindowWidth {
		width = minWindowWidth
	}

	if height < minWindowHeight {
		height = minWindowHeight
	}

	win.Width = width
	win.Height = height
}

// sortedWindows is the lock-free internal version of Windows.
func (m *Manager) sortedWindows() []*Window {
	sorted := make([]*Window, len(m.windows))
	copy(sorted, m.windows)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].ZIndex < sorted[j].ZIndex
	})

	return sorted
}

// cycleFocus is the lock-free internal version of CycleFocus.
func (m *Manager) cycleFocus() {
	sorted := m.sortedWindows()
	if len(sorted) < 2 {
		return
	}

	for i, win := range sorted {
		if win.Focused {
			next := sorted[(i+1)%len(sorted)]
			m.focus(next.ID)

			return
		}
	}

	m.focus(sorted[0].ID)
}

// renderWindow draws a window frame with title bar and content.
func (m *Manager) renderWindow(win *Window) string {
	// Inner content area (subtract 2 for borders, 1 for title bar).
	innerW := win.Width - 2
	innerH := win.Height - 3

	if innerW < 1 {
		innerW = 1
	}

	if innerH < 1 {
		innerH = 1
	}

	// Render the tea.Model content.
	content := win.Content.View()

	// Build title bar.
	titleBar := m.buildTitleBar(win, innerW)

	// Build resize grip.
	grip := m.buildResizeGrip(innerW)

	// Frame the content with a border — blue background fills the side borders.
	bgColor := config.Values.DialogBoxStyle.GetBackground()
	border := config.Values.DialogBoxStyle.GetBorderStyle()
	borderStyle := lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Background(bgColor)
	styledLeft := borderStyle.Render(border.Left)
	styledRight := borderStyle.Render(border.Right)

	framedContent := lipgloss.NewStyle().
		Width(innerW).
		Height(innerH).
		Background(bgColor).
		Render(content)

	// Compose: title bar + content + grip, with side borders.
	var out strings.Builder

	out.WriteString(titleBar)
	out.WriteByte('\n')

	// Content lines with side borders.
	padStyle := lipgloss.NewStyle().Background(bgColor)

	for line := range strings.SplitSeq(framedContent, "\n") {
		out.WriteString(styledLeft)

		// Pad or clip line to innerW.
		lineW := lipgloss.Width(line)
		if lineW < innerW {
			line += padStyle.Render(strings.Repeat(" ", innerW-lineW))
		}

		out.WriteString(line)
		out.WriteString(styledRight)
		out.WriteByte('\n')
	}

	// Bottom border with resize grip.
	out.WriteString(grip)

	return out.String()
}

func (m *Manager) buildTitleBar(win *Window, innerW int) string {
	border := config.Values.DialogBoxStyle.GetBorderStyle()
	bgColor := config.Values.DialogBoxStyle.GetBackground()
	borderFg := config.Values.DialogBoxStyle.GetBorderBottomForeground()

	// NC-style title: [ Title ]
	title := "[ " + win.Title + " ]"

	titleLen := lipgloss.Width(title)
	if titleLen > innerW {
		title = title[:innerW]
		titleLen = innerW
	}

	padding := innerW - titleLen
	left := padding / 2
	right := padding - left

	borderStyle := lipgloss.NewStyle().Foreground(borderFg).Background(bgColor)

	var titleStyle lipgloss.Style
	if win.Focused {
		// Focused: bright white title, bold — matches NC active dialog style.
		titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(bgColor).
			Bold(true)
	} else {
		// Unfocused: same cyan as the border, no bold.
		titleStyle = lipgloss.NewStyle().
			Foreground(borderFg).
			Background(bgColor)
	}

	leftPart := border.TopLeft + strings.Repeat(border.Top, left)
	rightPart := strings.Repeat(border.Top, right) + border.TopRight

	return borderStyle.Render(leftPart) + titleStyle.Render(title) + borderStyle.Render(rightPart)
}

func (m *Manager) buildResizeGrip(innerW int) string {
	border := config.Values.DialogBoxStyle.GetBorderStyle()
	bgColor := config.Values.DialogBoxStyle.GetBackground()

	gripRune := "◢"
	bottomLen := max(innerW-1, 0)

	grip := border.BottomLeft +
		strings.Repeat(border.Bottom, bottomLen) +
		gripRune +
		border.BottomRight

	return lipgloss.NewStyle().
		Foreground(config.Values.DialogBoxStyle.GetBorderBottomForeground()).
		Background(bgColor).
		Render(grip)
}

func (m *Manager) handleMouse(msg tea.MouseMsg) tea.Cmd {
	mouseX, mouseY := msg.X, msg.Y

	switch msg.Action {
	case tea.MouseActionPress:
		if msg.Button != tea.MouseButtonLeft {
			return nil
		}

		return m.handleMousePress(mouseX, mouseY, msg)

	case tea.MouseActionMotion:
		for _, win := range m.windows {
			if win.dragging {
				win.X = mouseX - win.dragOffX
				win.Y = mouseY - win.dragOffY

				return nil
			}

			if win.resizing {
				deltaW := mouseX - win.resizeOrigX
				deltaH := mouseY - win.resizeOrigY
				m.resize(win.ID, win.resizeOrigW+deltaW, win.resizeOrigH+deltaH)

				return nil
			}
		}

	case tea.MouseActionRelease:
		for _, win := range m.windows {
			win.dragging = false
			win.resizing = false
		}
	}

	return nil
}

func (m *Manager) handleMousePress(mouseX, mouseY int, msg tea.MouseMsg) tea.Cmd {
	sorted := m.sortedWindows()

	for i := range slices.Backward(sorted) {
		win := sorted[i]
		if !win.Visible || !win.Contains(mouseX, mouseY) {
			continue
		}

		m.focus(win.ID)

		if win.InTitleBar(mouseX, mouseY) {
			win.dragging = true
			win.dragOffX = mouseX - win.X
			win.dragOffY = mouseY - win.Y

			return nil
		}

		if win.InResizeGrip(mouseX, mouseY) {
			win.resizing = true
			win.resizeOrigW = win.Width
			win.resizeOrigH = win.Height
			win.resizeOrigX = mouseX
			win.resizeOrigY = mouseY

			return nil
		}

		// Click inside content area — translate coordinates and forward.
		localMsg := tea.MouseMsg{
			X:      mouseX - win.X - 1,
			Y:      mouseY - win.Y - 1,
			Action: msg.Action,
			Button: msg.Button,
		}

		var cmd tea.Cmd

		win.Content, cmd = win.Content.Update(localMsg)

		return cmd
	}

	return nil
}

func (m *Manager) handleKey(msg tea.KeyMsg) tea.Cmd {
	// Route key events only to the focused window.
	for _, win := range m.windows {
		if win.Focused && win.Visible {
			var cmd tea.Cmd

			win.Content, cmd = win.Content.Update(msg)

			return cmd
		}
	}

	return nil
}

// Helpers.

func (m *Manager) get(winID int) *Window {
	for _, win := range m.windows {
		if win.ID == winID {
			return win
		}
	}

	return nil
}

func (m *Manager) topZIndex() int {
	zIndex := 0

	for _, win := range m.windows {
		if win.ZIndex > zIndex {
			zIndex = win.ZIndex
		}
	}

	return zIndex
}

func (m *Manager) topWindow() *Window {
	var top *Window

	for _, win := range m.windows {
		if top == nil || win.ZIndex > top.ZIndex {
			top = win
		}
	}

	return top
}
