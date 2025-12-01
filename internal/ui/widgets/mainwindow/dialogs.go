package mainwindow

import (
	"context"

	"github.com/drekunov/gc/internal/ui/widgets/info"
)

func (m *Model) ShowModal(title, message string) {
	//TODO implement me
	panic("implement me")
}

func (m *Model) Info(title, footer, message string) {
	m.about = info.New()
	m.about.SetText(message)
	m.about.SetTitle(title)
	m.about.SetFooter(footer)
	m.about.SetVisible(true)
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

func (m *Model) Input(ctx context.Context, title, footer, message string) string {
	m.about = info.New()
	m.about.SetText(message)
	m.about.SetTitle(title)
	m.about.SetFooter(footer)
	m.about.SetVisible(true)

	select {
	case <-ctx.Done():
	case <-m.about.Done():
	}

	return m.about.Input()
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
