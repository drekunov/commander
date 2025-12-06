package mainwindow

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	windowID    int64
	windowsList map[windowID]tea.Model
)

func (w windowsList) Init() tea.Cmd {
	return nil
}

func (w windowsList) Update(msg tea.Msg) (windowsList, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	for id, model := range w {
		w[id], cmd = model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return w, tea.Batch(cmds...)
}

func (w windowsList) View() string {
	var views []string

	for _, model := range w {
		views = append(views, model.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, views...)
}

func (m *Model) addWindow(w tea.Model) windowID {
	wid := newWindowID()
	m.windows[wid] = w

	return wid
}

func (m *Model) removeWindow(id windowID) {
	delete(m.windows, id)
}

func newWindowID() windowID {
	return windowID(time.Now().UnixNano())
}
