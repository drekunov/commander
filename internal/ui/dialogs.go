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

	id := m.main.AddWindow(about)
	defer m.main.RemoveWindow(id)

	select {
	case <-ctx.Done():
	case <-about.Done():
	}

	m.program.Send(tea.ResumeMsg{})
}

func (m *Model) Input(ctx context.Context, title, footer, message string) string {
	input := dialogs.NewInput()
	input.SetText(message)
	input.SetTitle(title)
	input.SetFooter(footer)
	input.SetVisible(true)

	id := m.main.AddWindow(input)
	defer m.main.RemoveWindow(id)

	m.program.Send(tea.ResumeMsg{})

	select {
	case <-ctx.Done():
	case <-input.Done():
	}

	return input.Input()
}

func (m *Model) Password(title, message string) string {
	// TODO implement me
	panic("implement me")
}

func (m *Model) Select(title, message string, options []string) string {
	// TODO implement me
	panic("implement me")
}

func (m *Model) SelectMultiple(title, message string, options []string) []string {
	// TODO implement me
	panic("implement me")
}

func (m *Model) Error(title, message string) {
	// TODO implement me
	panic("implement me")
}

func (m *Model) Warning(title, message string) {
	// TODO implement me
	panic("implement me")
}

func (m *Model) Confirm(title, message string) bool {
	// TODO implement me
	panic("implement me")
}
