//nolint:testpackage // tests read the dialog's unexported visible state
package dialogs

import (
	"io"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/muesli/termenv"
)

func TestSelectEscapeCancels(t *testing.T) {
	t.Parallel()

	sel := NewSelect(config.Styles{}, []string{"Name", "Size"}, false)
	sel.SetVisible(true)

	_, _ = sel.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if !sel.Canceled() {
		t.Error("Escape did not mark the dialog canceled")
	}

	if sel.visible {
		t.Error("Escape did not hide the dialog")
	}

	select {
	case <-sel.Done():
	default:
		t.Error("Escape did not signal Done")
	}
}

func TestSelectCursorUsesCursorStyle(t *testing.T) {
	t.Parallel()

	renderer := lipgloss.NewRenderer(io.Discard)
	renderer.SetColorProfile(termenv.ANSI256)

	cursor := lipgloss.NewStyle().Underline(true).Renderer(renderer)
	sel := NewSelect(config.Styles{CursorStyle: cursor}, []string{"Name", "Size"}, false)
	sel.SetVisible(true)

	view := sel.View()

	if !strings.Contains(view, "\x1b[4m") {
		t.Errorf("focused row is not styled with the cursor style:\n%q", view)
	}
}
