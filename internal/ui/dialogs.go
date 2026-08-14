package ui

import (
	"context"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
)

const (
	dialogDefaultWidth  = 40
	dialogDefaultHeight = 20
)

func (m *Model) Info(ctx context.Context, title, footer, message string) {
	about := dialogs.NewInfo(m.styles)
	about.SetText(message)
	about.SetTitle(title)
	about.SetFooter(footer)
	about.SetVisible(true)

	id := m.addDialogWindow(about, title)
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
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
	input.SetTitle(title)
	input.SetFooter(footer)
	input.SetVisible(true)

	id := m.addDialogWindow(input, title)
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return ""
	case <-input.Done():
		return input.Input()
	}
}

func (m *Model) Select(ctx context.Context, title, message string, options []string) string {
	sel := dialogs.NewSelect(m.styles, options, false)
	sel.SetText(message)
	sel.SetTitle(title)
	sel.SetVisible(true)

	id := m.addDialogWindow(sel, title)
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return ""
	case <-sel.Done():
		return sel.Selection()
	}
}

func (m *Model) SelectMultiple(ctx context.Context, title, message string, options []string) []string {
	sel := dialogs.NewSelect(m.styles, options, true)
	sel.SetText(message)
	sel.SetTitle(title)
	sel.SetVisible(true)

	id := m.addDialogWindow(sel, title)
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
		return nil
	case <-sel.Done():
		return sel.Selections()
	}
}

func (m *Model) Confirm(ctx context.Context, title, message string) bool {
	confirm := dialogs.NewConfirm(m.styles, message)
	confirm.SetTitle(title)
	confirm.SetVisible(true)

	id := m.addDialogWindow(confirm, title)
	defer m.main.WM().Remove(id)

	m.sendMsg(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
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

// addDialogWindow adds a tea.Model as a centered dialog window via the wm.
func (m *Model) addDialogWindow(content tea.Model, title string) int {
	wmgr := m.main.WM()
	width := wmgr.Width()
	height := wmgr.Height()

	if width < 2 {
		width = dialogDefaultWidth
	}

	if height < 2 {
		height = dialogDefaultHeight
	}

	winW := width / 2
	winH := height / 2
	posX := (width - winW) / 2
	posY := (height - winH) / 2

	return wmgr.Add(content, title, posX, posY, winW, winH)
}
