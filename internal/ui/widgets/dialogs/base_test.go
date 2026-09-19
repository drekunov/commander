//nolint:testpackage // tests exercise the unexported dialog base
package dialogs

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

func TestDialogRenderHasNoBorderAndKeepsFooter(t *testing.T) {
	t.Parallel()

	// Inject a bordered dialog style so the test would catch the base drawing
	// its own frame; with a zero style the assertion would pass vacuously.
	bordered := config.Styles{
		DialogBoxStyle: lipgloss.NewStyle().Border(lipgloss.DoubleBorder()),
	}

	base := newDialogBase(bordered)
	base.SetFooter("[Error]")

	out := base.render("body line")

	if strings.ContainsAny(out, "╔╗╚╝║═") {
		t.Errorf("render output contains border characters:\n%s", out)
	}

	if !strings.Contains(out, "body line") {
		t.Errorf("render output is missing the body:\n%s", out)
	}

	if !strings.Contains(out, "[Error]") {
		t.Errorf("render output is missing the footer:\n%s", out)
	}
}
