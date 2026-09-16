//nolint:testpackage // tests construct Model directly to inject the navigator
package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

func TestNavigateMsgInvokesNavigator(t *testing.T) {
	t.Parallel()

	model := &Model{
		main:      mainform.New(config.Styles{}),
		styles:    config.Styles{},
		sendMsg:   func(tea.Msg) {},
		navigator: nil,
	}

	type nav struct {
		panel app.PanelID
		dir   string
	}

	got := make(chan nav, 1)

	model.navigator = func(panel app.PanelID, dir string) {
		got <- nav{panel: panel, dir: dir}
	}

	_, cmd := model.Update(panel.NavigateMsg{Panel: app.PanelRight, Dir: "/tmp"})
	if cmd != nil {
		t.Errorf("NavigateMsg produced a command %v, want none", cmd())
	}

	select {
	case nav := <-got:
		if nav.panel != app.PanelRight {
			t.Errorf("navigator received panel %v, want right", nav.panel)
		}

		if nav.dir != "/tmp" {
			t.Errorf("navigator received dir %q, want /tmp", nav.dir)
		}
	default:
		t.Fatal("navigator was not invoked")
	}
}
