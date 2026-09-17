//nolint:testpackage // tests construct Model directly to reach the focused panel
package ui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/app"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

const (
	sortAttrName  = "Name"
	sortAttrSize  = "Size"
	sortAttrIsDir = "IsDir"
)

// sortModel builds a Model whose focused (left) panel has the four columns.
func sortModel(t *testing.T) *Model {
	t.Helper()

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
	}

	header := app.AttributeList{
		{AttrName: sortAttrName, AttrValue: sortAttrName, Width: 11, Flex: true},
		{AttrName: sortAttrSize, AttrValue: sortAttrSize, Width: 6},
		{AttrName: sortAttrIsDir, AttrValue: sortAttrIsDir, Hidden: true},
	}

	model.main.Panel().SetData("/", []app.AttributeList{header})

	return model
}

func TestMenuActivationOpensSortWindow(t *testing.T) {
	t.Parallel()

	model := sortModel(t)

	if cmd := model.handleBarActivation(buttonbar.ActionMenu); cmd == nil {
		t.Fatal("Menu produced no command, want the sort window command")
	}

	if model.mockDialog != nil {
		t.Error("Menu activation opened a mock dialog")
	}
}

func TestSortColumnMsgAppliesSort(t *testing.T) {
	t.Parallel()

	model := sortModel(t)

	if col, asc := model.main.Panel().SortState(); col != 0 || !asc {
		t.Fatalf("initial sort = %d/%v, want Name ascending", col, asc)
	}

	model.Update(sortColumnMsg{column: sortAttrSize})

	if col, asc := model.main.Panel().SortState(); col != 1 || !asc {
		t.Errorf("sort = %d/%v, want Size ascending", col, asc)
	}
}

func TestSortColumnMsgIgnoresCancel(t *testing.T) {
	t.Parallel()

	model := sortModel(t)

	model.Update(sortColumnMsg{column: ""})

	if col, asc := model.main.Panel().SortState(); col != 0 || !asc {
		t.Errorf("sort = %d/%v, want the unchanged Name ascending", col, asc)
	}
}

func TestSortWindowCaptionCarriesPrompt(t *testing.T) {
	t.Parallel()

	model := sortModel(t)

	cmd := model.sortWindowCmd()
	if cmd == nil {
		t.Fatal("sortWindowCmd returned nil")
	}

	done := make(chan struct{})

	go func() {
		cmd()
		close(done)
	}()

	win := waitForSelectWindow(model)
	if win == nil {
		t.Fatal("the sort window did not open")
	}

	const wantCaption = "Sort by"

	if win.Title != wantCaption {
		t.Errorf("caption = %q, want %q", win.Title, wantCaption)
	}

	sel, ok := win.Content.(*dialogs.Select)
	if !ok {
		t.Fatalf("window content = %T, want *dialogs.Select", win.Content)
	}

	if strings.Contains(sel.View(), "Select a column") {
		t.Errorf("the body repeats the prompt:\n%s", sel.View())
	}

	// Close the dialog so the command goroutine returns.
	_, _ = sel.Update(tea.KeyMsg{Type: tea.KeyEnter})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("the sort window command did not return")
	}
}

// waitForSelectWindow polls the window manager until a Select dialog window
// appears, up to a short timeout.
func waitForSelectWindow(model *Model) *wm.Window {
	for range 200 {
		for _, win := range model.main.WM().Windows() {
			if _, ok := win.Content.(*dialogs.Select); ok {
				return win
			}
		}

		time.Sleep(time.Millisecond)
	}

	return nil
}
