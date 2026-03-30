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
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		n.handleKey(msg)
	}

	return n, nil
}

func (n *notepad) View() string {
	var output strings.Builder

	for i, line := range n.lines {
		if i == n.cur {
			output.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFF00")).
				Render("> " + line + "_"))
		} else {
			output.WriteString("  " + line)
		}

		output.WriteByte('\n')
	}

	return output.String()
}

func (n *notepad) handleKey(msg tea.KeyMsg) {
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
		n.handleBackspace()
	case tea.KeyRunes:
		n.lines[n.cur] += string(msg.Runes)
	}
}

func (n *notepad) handleBackspace() {
	line := n.lines[n.cur]
	if len(line) > 0 {
		n.lines[n.cur] = line[:len(line)-1]
	} else if n.cur > 0 {
		n.lines = append(n.lines[:n.cur], n.lines[n.cur+1:]...)
		n.cur--
	}
}

// listView is a scrollable list widget.
type listView struct {
	items    []string
	selected int
}

func newListView(items []string) *listView {
	return &listView{items: items}
}

func (lv *listView) Init() tea.Cmd { return nil }

func (lv *listView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) { //nolint: gocritic
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if lv.selected > 0 {
				lv.selected--
			}
		case tea.KeyDown:
			if lv.selected < len(lv.items)-1 {
				lv.selected++
			}
		}
	}

	return lv, nil
}

func (lv *listView) View() string {
	var output strings.Builder

	for i, item := range lv.items {
		if i == lv.selected {
			output.WriteString(lipgloss.NewStyle().
				Foreground(lipgloss.Color("#000000")).
				Background(lipgloss.Color("#00AAAA")).
				Render(" " + item + " "))
		} else {
			output.WriteString("  " + item)
		}

		output.WriteByte('\n')
	}

	return output.String()
}

// infoWidget displays a static message.
type infoWidget struct {
	message string
}

func (iw *infoWidget) Init() tea.Cmd { return nil }

func (iw *infoWidget) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return iw, nil
}

func (iw *infoWidget) View() string {
	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Bold(true)

	return style.Render(iw.message)
}

// --- Main application model ---

type model struct {
	manager *wm.Manager
}

func newModel() *model {
	app := &model{
		manager: wm.New(),
	}

	// Window 1: Notepad at top-left.
	app.manager.Add(
		newNotepad("notes", "Hello from the window manager!\nTry clicking other windows.\nType here to edit."),
		"Notepad",
		2, 1, 40, 15,
	)

	// Window 2: List view, overlapping the first.
	app.manager.Add(
		newListView([]string{
			"Documents", "Downloads", "Pictures",
			"Music", "Videos", "Desktop",
			"Projects", "Trash",
		}),
		"File List",
		20, 5, 30, 14,
	)

	helpText := "Window Manager Example\n\n" +
		"Drag title bars to move.\n" +
		"Drag ◢ to resize.\n" +
		"Click to focus.\n" +
		"Tab to cycle windows.\n" +
		"F10 to quit."

	// Window 3: Small info widget at the right.
	app.manager.Add(
		&infoWidget{message: helpText},
		"Help",
		50, 2, 32, 12,
	)

	return app
}

func (app *model) Init() tea.Cmd {
	return app.manager.Init()
}

func (app *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyF10:
			return app, tea.Quit
		case tea.KeyTab:
			app.cycleWindows()

			return app, nil
		}
	case tea.WindowSizeMsg:
		app.manager.SetSize(msg.Width, msg.Height)
	}

	cmd := app.manager.Update(msg)

	return app, cmd
}

func (app *model) View() string {
	return app.manager.View()
}

func (app *model) cycleWindows() {
	windows := app.manager.Windows()
	if len(windows) < 2 {
		return
	}

	// Find focused, focus the next one.
	for i, win := range windows {
		if win.Focused {
			next := windows[(i+1)%len(windows)]
			app.manager.Focus(next.ID)

			return
		}
	}

	app.manager.Focus(windows[0].ID)
}

func main() {
	program := tea.NewProgram(
		newModel(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	)

	_, err := program.Run()
	if err != nil {
		log.Fatal(fmt.Errorf("error running program: %w", err))
	}
}
