package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
)

func (m *Model) Info(ctx context.Context, title, footer, message string) {
	about := dialogs.NewInfo()
	about.SetText(message)
	about.SetTitle(title)
	about.SetFooter(footer)
	about.SetVisible(true)

	id := m.addDialogWindow(about, title)
	defer m.main.WM().Remove(id)

	m.program.Send(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
	case <-about.Done():
	}
}

func (m *Model) Input(ctx context.Context, title, footer, message string) string {
	input := dialogs.NewInput()
	input.SetText(message)
	input.SetTitle(title)
	input.SetFooter(footer)
	input.SetVisible(true)

	id := m.addDialogWindow(input, title)
	defer m.main.WM().Remove(id)

	m.program.Send(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
	case <-input.Done():
	}

	return input.Input()
}

func (m *Model) Password(title, message string) string {
	return m.Input(context.Background(), title, "", message)
}

func (m *Model) Select(title, message string, Koptions []string) string {
	m.Info(context.Background(), title, "", message)

	return ""
}

func (m *Model) SelectMultiple(title, message string, options []string) []string {
	m.Info(context.Background(), title, "", message)

	return nil
}

func (m *Model) Error(title, message string) {
	m.Info(context.Background(), title, "[Error]", message)
}

func (m *Model) Warning(title, message string) {
	m.Info(context.Background(), title, "[Warning]", message)
}

func (m *Model) Confirm(title, message string) bool {
	input := dialogs.NewInput()
	input.SetText(message)
	input.SetTitle(title)
	input.SetFooter("Y/n")
	input.SetVisible(true)

	id := m.addDialogWindow(input, title)
	defer m.main.WM().Remove(id)

	m.program.Send(tea.ResumeMsg{})

	<-input.Done()

	return input.Input() == "Y"
}

// addDialogWindow adds a tea.Model as a centered dialog window via the wm.
func (m *Model) addDialogWindow(content tea.Model, title string) int {
	wmgr := m.main.WM()
	width := wmgr.Width()
	height := wmgr.Height()

	if width < 2 {
		width = 40
	}

	if height < 2 {
		height = 20
	}

	winW := width / 2
	winH := height / 2
	posX := (width - winW) / 2
	posY := (height - winH) / 2

	return wmgr.Add(content, title, posX, posY, winW, winH)
}
