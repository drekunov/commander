package ui

import tea "github.com/charmbracelet/bubbletea"

func (m *Model) ShowModal(title, message string) {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Info(title, message string) {
	m.main.Info(title, message)
	m.program.Send(tea.ResumeMsg{})
}

func (m *Model) Error(title, message string) {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Warning(title, message string) {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Confirm(title, message string) bool {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Input(title, message string) string {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Password(title, message string) string {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Select(title, message string, options []string) string {
	//TODO implement me
	panic("implement me")
}

func (m *Model) SelectMultiple(title, message string, options []string) []string {
	//TODO implement me
	panic("implement me")
}
