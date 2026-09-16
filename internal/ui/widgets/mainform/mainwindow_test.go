package mainform_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/panel"
)

func resizeMainform(m *mainform.Model) *mainform.Model {
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	return updated
}

func TestFunctionKeyReturnsActivateMsg(t *testing.T) {
	t.Parallel()

	model := mainform.New(config.Styles{})
	model = resizeMainform(model)

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyF5})
	if cmd == nil {
		t.Fatal("F5 produced no command")
	}

	msg, ok := cmd().(buttonbar.ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want buttonbar.ActivateMsg", cmd())
	}

	if msg.Action != buttonbar.ActionCopy {
		t.Errorf("F5 resolved to %v, want ActionCopy", msg.Action)
	}
}

func TestBarRowPressIsConsumed(t *testing.T) {
	t.Parallel()

	model := mainform.New(config.Styles{})
	model = resizeMainform(model)

	updated, cmd := model.Update(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      44,
		Y:      29,
	})

	if cmd == nil {
		t.Fatal("bar-row press produced no command")
	}

	if !updated.Bar().Focused() {
		t.Error("bar did not gain focus after a bar-row press")
	}

	msg, ok := cmd().(buttonbar.ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want buttonbar.ActivateMsg", cmd())
	}

	if msg.Action != buttonbar.ActionCopy {
		t.Errorf("bar-row press at column 44 resolved to %v, want ActionCopy", msg.Action)
	}
}

func TestPressAboveBarReachesWmAndClearsBarFocus(t *testing.T) {
	t.Parallel()

	model := mainform.New(config.Styles{})
	model = resizeMainform(model)

	model.Bar().SetFocused(true)

	updated, cmd := model.Update(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      55,
		Y:      5,
	})

	if cmd != nil {
		t.Errorf("above-bar press produced a command %v, want none", cmd())
	}

	if updated.Bar().Focused() {
		t.Error("above-bar press did not clear bar focus")
	}

	left := updated.Panel()
	focusLeft, focusOther := false, false

	for _, win := range updated.WM().Windows() {
		if !win.Focused {
			continue
		}

		if win.Content == left {
			focusLeft = true
		} else {
			focusOther = true
		}
	}

	if !focusOther {
		t.Error("above-bar press on the right panel did not deliver focus to a window")
	}

	if focusLeft {
		t.Error("above-bar press left the left panel focused instead of the clicked window")
	}
}

func listingFor(names ...string) []app.AttributeList {
	data := []app.AttributeList{
		{
			{AttrName: "Name", AttrValue: "Name"},
			{AttrName: "IsDir", AttrValue: "IsDir"},
		},
	}

	for _, name := range names {
		data = append(data, app.AttributeList{
			{AttrName: "Name", AttrValue: name},
			{AttrName: "IsDir", AttrValue: false},
		})
	}

	return data
}

func TestDataMsgRoutesToAddressedPanel(t *testing.T) {
	t.Parallel()

	model := mainform.New(config.Styles{})
	model = resizeMainform(model)

	_, cmd := model.Update(panel.DataMsg{
		Panel: app.PanelRight,
		Path:  "/etc",
		Data:  listingFor("hosts"),
	})

	if cmd != nil {
		t.Errorf("DataMsg produced a command %v, want none", cmd())
	}

	rightUpdated, leftUpdated := false, false

	for _, win := range model.WM().Windows() {
		content, ok := win.Content.(*panel.Model)
		if !ok {
			continue
		}

		switch content.Dir() {
		case "/etc":
			rightUpdated = true

			if win.Title != "/etc" {
				t.Errorf("right window title = %q, want /etc", win.Title)
			}
		case "":
			leftUpdated = true
		}
	}

	if !rightUpdated {
		t.Error("right panel was not updated to /etc")
	}

	if !leftUpdated {
		t.Error("left panel was unexpectedly changed")
	}
}
