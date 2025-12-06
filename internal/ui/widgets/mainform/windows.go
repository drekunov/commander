package mainform

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	WindowID    int64
	windowsList map[WindowID]tea.Model
)

func (w windowsList) Init() tea.Cmd {
	return nil
}

func (w windowsList) Update(msg tea.Msg) (windowsList, tea.Cmd) {
	var cmd tea.Cmd

	cmds := make([]tea.Cmd, 0, len(w))

	for id, model := range w {
		w[id], cmd = model.Update(msg)
		cmds = append(cmds, cmd)
	}

	return w, tea.Batch(cmds...)
}

func (w windowsList) View() string {
	views := make([]string, 0, len(w))

	for _, model := range w {
		views = append(views, model.View())
	}

	return lipgloss.JoinVertical(lipgloss.Center, views...)
}

func (m *Model) AddWindow(w tea.Model) WindowID {
	wid := newWindowID()
	m.windows[wid] = w

	return wid
}

func (m *Model) RemoveWindow(id WindowID) {
	delete(m.windows, id)
}

func newWindowID() WindowID {
	return WindowID(time.Now().UnixNano())
}
