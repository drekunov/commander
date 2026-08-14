package wm

import (
	"slices"

	tea "github.com/charmbracelet/bubbletea"
)

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
