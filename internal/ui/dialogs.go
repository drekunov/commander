package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Info(ctx context.Context, title, footer, message string) {
	m.main.Info(ctx, title, footer, message)

	m.program.Send(tea.ResumeMsg{})
}

func (m *Model) Input(ctx context.Context, title, footer, message string) string {
	return m.main.Input(ctx, title, footer, message)
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
