package main

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

// --- Example tea.Model widgets to put inside windows ---

// notepad is a simple text editor widget.
type notepad struct {
	title string
	lines []string
	cur   int
}

func newNotepad(title string, text string) *notepad {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	return &notepad{title: title, lines: lines}
}

func (n *notepad) Init() tea.Cmd { return nil }

func (n *notepad) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if n.cur > 0 {
				n.cur--
			}
		case tea.KeyDown:
			if n.cur < len(n.lines)-1 {
				n.cur++
			}
		case tea.KeyEnter:
			n.lines = append(n.lines[:n.cur+1], append([]string{""}, n.lines[n.cur+1:]...)...)
			n.cur++
		case tea.KeyBackspace:
			line := n.lines[n.cur]
			if len(line) > 0 {
				n.lines[n.cur] = line[:len(line)-1]
			} else if n.cur > 0 {
				n.lines = append(n.lines[:n.cur], n.lines[n.cur+1:]...)
				n.cur--
			}
		default:
			if msg.Type == tea.KeyRunes {
				n.lines[n.cur] += string(msg.Runes)
			}
		}
	}
	return n, nil
}

func (n *notepad) View() string {
	var b strings.Builder
	for i, line := range n.lines {
		if i == n.cur {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFF00")).
				Render("> "+line+"_"))
		} else {
			b.WriteString("  " + line)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// listView is a scrollable list widget.
type listView struct {
	items    []string
	selected int
}

func newListView(items []string) *listView {
	return &listView{items: items}
}

func (l *listView) Init() tea.Cmd { return nil }

func (l *listView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if l.selected > 0 {
				l.selected--
			}
		case tea.KeyDown:
			if l.selected < len(l.items)-1 {
				l.selected++
			}
		}
	}
	return l, nil
}

func (l *listView) View() string {
	var b strings.Builder
	for i, item := range l.items {
		if i == l.selected {
			b.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#00AAAA")).
				Render(" "+item+" "))
		} else {
			b.WriteString("  " + item)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

// clockWidget displays a static message (placeholder for any widget).
type clockWidget struct {
	message string
}

func (c *clockWidget) Init() tea.Cmd { return nil }

func (c *clockWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return c, nil
}

func (c *clockWidget) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Bold(true)
	return style.Render(c.message)
}

// --- Main application model ---

type model struct {
	manager *wm.Manager
}

func newModel() *model {
	m := &model{
		manager: wm.New(),
	}

	// Window 1: Notepad at top-left.
	m.manager.Add(
		newNotepad("notes", "Hello from the window manager!\nTry clicking other windows.\nType here to edit."),
		"Notepad",
		2, 1, 40, 15,
	)

	// Window 2: List view, overlapping the first.
	m.manager.Add(
		newListView([]string{
			"Documents", "Downloads", "Pictures",
			"Music", "Videos", "Desktop",
			"Projects", "Trash",
		}),
		"File List",
		20, 5, 30, 14,
	)

	// Window 3: Small info widget at the right.
	m.manager.Add(
		&clockWidget{message: "Window Manager Example\n\nDrag title bars to move.\nDrag ◢ to resize.\nClick to focus.\nTab to cycle windows.\nF10 to quit."},
		"Help",
		50, 2, 32, 12,
	)

	return m
}

func (m *model) Init() tea.Cmd {
	return m.manager.Init()
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF10:
			return m, tea.Quit
		case tea.KeyTab:
			m.cycleWindows()
			return m, nil
		}
	case tea.WindowSizeMsg:
		m.manager.SetSize(msg.Width, msg.Height)
	}

	cmd := m.manager.Update(msg)
	return m, cmd
}

func (m *model) View() string {
	return m.manager.View()
}

func (m *model) cycleWindows() {
	windows := m.manager.Windows()
	if len(windows) < 2 {
		return
	}
	// Find focused, focus the next one.
	for i, w := range windows {
		if w.Focused {
			next := windows[(i+1)%len(windows)]
			m.manager.Focus(next.ID)
			return
		}
	}
	m.manager.Focus(windows[0].ID)
}

func main() {
	p := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatal(fmt.Errorf("error running program: %w", err))
	}
}
