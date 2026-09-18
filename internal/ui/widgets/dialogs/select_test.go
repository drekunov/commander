//nolint:testpackage // tests read the dialog's unexported visible state
package dialogs

import (
	"io"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
	"github.com/muesli/termenv"
)

const (
	optionName = "Name"
	optionSize = "Size"
)

func TestSelectEscapeCancels(t *testing.T) {
	t.Parallel()

	sel := NewSelect(config.Styles{}, []string{optionName, optionSize}, false)
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

func TestSelectUsesDialogOptionStyles(t *testing.T) {
	t.Parallel()

	renderer := lipgloss.NewRenderer(io.Discard)
	renderer.SetColorProfile(termenv.ANSI256)

	dialogCursor := lipgloss.NewStyle().Underline(true).Renderer(renderer)
	dialogOption := lipgloss.NewStyle().Italic(true).Renderer(renderer)

	sel := NewSelect(config.Styles{
		DialogCursorStyle: dialogCursor,
		DialogOptionStyle: dialogOption,
	}, []string{optionName, optionSize}, false)
	sel.SetVisible(true)

	view := sel.View()

	if !strings.Contains(view, "\x1b[4") {
		t.Errorf("focused row is not styled with the dialog cursor style:\n%q", view)
	}

	if !strings.Contains(view, "\x1b[3") {
		t.Errorf("unfocused row is not styled with the dialog option style:\n%q", view)
	}
}

func TestSelectCursorFallsBackToPanelCursor(t *testing.T) {
	t.Parallel()

	cursor := lipgloss.NewStyle().Underline(true)
	sel := NewSelect(config.Styles{CursorStyle: cursor}, []string{optionName}, false)

	if got := sel.optionStyle(true); !got.GetUnderline() {
		t.Error("focused option did not fall back to the panel cursor style")
	}
}

func TestDialogBodyUsesFrameBackground(t *testing.T) {
	t.Parallel()

	background := lipgloss.Color("#000080")
	sel := NewSelect(config.Styles{
		DialogBoxStyle: lipgloss.NewStyle().Background(background),
	}, []string{optionName}, false)

	if got := sel.textStyle().GetBackground(); got != background {
		t.Errorf("dialog text background = %v, want %v", got, background)
	}

	if got := sel.optionStyle(true).GetBackground(); got != background {
		t.Errorf("dialog cursor background = %v, want %v", got, background)
	}

	if got := sel.optionStyle(false).GetBackground(); got != background {
		t.Errorf("dialog option background = %v, want %v", got, background)
	}
}

func TestSelectOptionsAreFlushLeftWithoutCursorMarker(t *testing.T) {
	t.Parallel()

	sel := NewSelect(config.Styles{}, []string{optionName, optionSize}, false)
	sel.SetVisible(true)

	lines := strings.Split(ansi.Strip(sel.View()), "\n")

	for _, line := range lines {
		if strings.Contains(line, ">") {
			t.Errorf("option line %q contains a cursor marker", line)
		}

		if strings.HasPrefix(line, " ") {
			t.Errorf("option line %q is not flush left", line)
		}
	}

	if lines[0] != optionName {
		t.Errorf("first option = %q, want %q flush left", lines[0], optionName)
	}
}

func TestSelectMultipleCheckboxIsFlushLeft(t *testing.T) {
	t.Parallel()

	sel := NewSelect(config.Styles{}, []string{optionName, optionSize}, true)
	sel.SetVisible(true)

	lines := strings.Split(ansi.Strip(sel.View()), "\n")

	if lines[0] != "[ ] "+optionName {
		t.Errorf("first option = %q, want %q", lines[0], "[ ] "+optionName)
	}
}

func TestSelectOptionsUseInjectedPadding(t *testing.T) {
	t.Parallel()

	option := lipgloss.NewStyle().Padding(0, 1)
	sel := NewSelect(config.Styles{
		DialogCursorStyle: option,
		DialogOptionStyle: option,
	}, []string{optionName, optionSize}, false)
	sel.SetVisible(true)

	lines := strings.Split(ansi.Strip(sel.View()), "\n")

	if lines[0] != " "+optionName+" " {
		t.Errorf("first option = %q, want %q", lines[0], " "+optionName+" ")
	}
}
