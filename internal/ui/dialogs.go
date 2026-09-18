package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

const (
	dialogDefaultWidth  = 40
	dialogDefaultHeight = 20
)

func (m *Model) Info(ctx context.Context, title, footer, message string) {
	about := dialogs.NewInfo(m.styles)
	about.SetText(message)
	about.SetFooter(footer)
	about.SetVisible(true)

	id := m.addDialogWindow(about, title)
	defer m.closeDialogWindow(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
	case <-m.quit:
	case <-about.Done():
	}
}

func (m *Model) Input(ctx context.Context, title, footer, message string) string {
	return m.runInput(ctx, title, footer, message, textinput.EchoNormal)
}

func (m *Model) Password(ctx context.Context, title, message string) string {
	return m.runInput(ctx, title, "", message, textinput.EchoPassword)
}

func (m *Model) runInput(ctx context.Context, title, footer, message string, echo textinput.EchoMode) string {
	input := dialogs.NewInput(m.styles, echo)
	input.SetText(message)
	input.SetFooter(footer)
	input.SetVisible(true)

	id := m.addDialogWindow(input, title)
	defer m.closeDialogWindow(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return ""
	case <-m.quit:
		return ""
	case <-input.Done():
		return input.Input()
	}
}

func (m *Model) Select(ctx context.Context, title, message string, options []string) string {
	sel := dialogs.NewSelect(m.styles, options, false)
	sel.SetText(message)
	sel.SetVisible(true)

	id := m.addDialogWindow(sel, title)
	defer m.closeDialogWindow(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return ""
	case <-m.quit:
		return ""
	case <-sel.Done():
		if sel.Canceled() {
			return ""
		}

		return sel.Selection()
	}
}

func (m *Model) SelectMultiple(ctx context.Context, title, message string, options []string) []string {
	sel := dialogs.NewSelect(m.styles, options, true)
	sel.SetText(message)
	sel.SetVisible(true)

	id := m.addDialogWindow(sel, title)
	defer m.closeDialogWindow(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return nil
	case <-m.quit:
		return nil
	case <-sel.Done():
		if sel.Canceled() {
			return nil
		}

		return sel.Selections()
	}
}

func (m *Model) Confirm(ctx context.Context, title, message string) bool {
	confirm := dialogs.NewConfirm(m.styles, message)
	confirm.SetVisible(true)

	id := m.addDialogWindow(confirm, title)
	defer m.closeDialogWindow(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return false
	case <-m.quit:
		return false
	case <-confirm.Done():
		return confirm.Decision()
	}
}

func (m *Model) Error(title, message string) {
	m.Info(context.Background(), title, "[Error]", message)
}

func (m *Model) Warning(title, message string) {
	m.Info(context.Background(), title, "[Warning]", message)
}

// addDialogWindow adds a tea.Model as a centered dialog window via the wm. The
// window is sized to fit the dialog content plus the window frame, so the frame
// hugs the content.
func (m *Model) addDialogWindow(content tea.Model, title string) int {
	wmgr := m.main.WM()

	canvasW := wmgr.Width()
	canvasH := wmgr.Height()

	if canvasW < 2 {
		canvasW = dialogDefaultWidth
	}

	if canvasH < 2 {
		canvasH = dialogDefaultHeight
	}

	view := content.View()
	innerW := max(lipgloss.Width(view), wm.TitleWidth(title))
	winW := innerW + wm.FrameCols
	winH := lipgloss.Height(view) + wm.FrameRows

	if winW > canvasW {
		winW = canvasW
	}

	if winH > canvasH {
		winH = canvasH
	}

	posX := (canvasW - winW) / 2
	posY := (canvasH - winH) / 2

	return wmgr.Add(content, title, posX, posY, winW, winH)
}

// closeDialogWindow removes a dialog window and wakes the event loop so the
// mainform re-syncs panel focus. Removing the window re-focuses the panel
// beneath it, but that happens off the event loop, so a message is needed to
// run the sync and restore the panel's cursor.
func (m *Model) closeDialogWindow(id int) {
	m.main.WM().Remove(id)
	m.sendMsg(tea.ResumeMsg{})
}
