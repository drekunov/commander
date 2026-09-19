//nolint:testpackage // tests construct Model directly to reach the dialog window
package ui

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
)

func TestSelectMultipleEscapeReturnsNil(t *testing.T) {
	t.Parallel()

	model := sortModel(t)
	result := make(chan []string, 1)

	go func() {
		result <- model.SelectMultiple(context.Background(), "Pick", "", []string{"a", "b"})
	}()

	win := waitForSelectWindow(model)
	if win == nil {
		t.Fatal("the select window did not open")
	}

	sel, ok := win.Content.(*dialogs.Select)
	if !ok {
		t.Fatalf("window content = %T, want *dialogs.Select", win.Content)
	}

	_, _ = sel.Update(tea.KeyMsg{Type: tea.KeyEsc})

	select {
	case got := <-result:
		if got != nil {
			t.Errorf("SelectMultiple after Escape = %v, want nil", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("SelectMultiple did not return after Escape")
	}
}
