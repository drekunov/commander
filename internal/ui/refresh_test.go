//nolint:testpackage // tests inject the refresher and read panel state
package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
)

func refreshTestModel(t *testing.T) *Model {
	t.Helper()

	return &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
	}
}

func TestCtrlRRefreshesBothPanels(t *testing.T) {
	t.Parallel()

	model := refreshTestModel(t)

	panels := model.main.Panels()
	if len(panels) != 2 {
		t.Fatalf("Panels() returned %d panels, want 2", len(panels))
	}

	panels[0].SetData("/left", nil)
	panels[1].SetData("/right", nil)

	var got []app.NavRequest

	model.SetRefresher(func(requests []app.NavRequest) {
		got = requests
	})

	model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})

	if len(got) != 2 {
		t.Fatalf("refresher got %d requests, want 2", len(got))
	}

	if got[0].Panel != app.PanelLeft || got[0].Dir != "/left" {
		t.Errorf("first request = %v/%q, want left /left", got[0].Panel, got[0].Dir)
	}

	if got[1].Panel != app.PanelRight || got[1].Dir != "/right" {
		t.Errorf("second request = %v/%q, want right /right", got[1].Panel, got[1].Dir)
	}
}

func TestCtrlRIgnoredWhileDialogOpen(t *testing.T) {
	t.Parallel()

	model := refreshTestModel(t)

	for _, p := range model.main.Panels() {
		p.SetData("/x", nil)
	}

	called := false

	model.SetRefresher(func([]app.NavRequest) {
		called = true
	})

	model.main.WM().Add(refreshDummyDialog{}, "Dialog", 1, 1, 10, 5)

	model.Update(tea.KeyMsg{Type: tea.KeyCtrlR})

	if called {
		t.Error("Ctrl+R invoked the refresher while a dialog was open")
	}
}

// refreshDummyDialog is a do-nothing window content used to stand in for a
// modal dialog.
type refreshDummyDialog struct{}

func (refreshDummyDialog) Init() tea.Cmd                       { return nil }
func (refreshDummyDialog) Update(tea.Msg) (tea.Model, tea.Cmd) { return refreshDummyDialog{}, nil }
func (refreshDummyDialog) View() string                        { return "" }
