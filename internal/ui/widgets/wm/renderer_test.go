//nolint:testpackage // tests exercise the unexported renderer and canvas
package wm

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
)

type staticModel struct {
	view string
}

func (s staticModel) Init() tea.Cmd {
	return nil
}

func (s staticModel) Update(_ tea.Msg) (tea.Model, tea.Cmd) {
	return s, nil
}

func (s staticModel) View() string {
	return s.view
}

func TestBuildTitleBar(t *testing.T) {
	t.Parallel()

	renderer := renderer{styles: config.Styles{}}

	win := NewWindow(1, staticModel{view: ""}, 0, 0, 20, 10)
	win.Title = "Files"
	win.Focused = true

	out := renderer.buildTitleBar(win, 18)

	if !strings.Contains(out, "[ Files ]") {
		t.Errorf("title bar %q missing the window title", out)
	}
}

func TestRenderWindowIncludesContent(t *testing.T) {
	t.Parallel()

	renderer := renderer{styles: config.Styles{}}

	win := NewWindow(1, staticModel{view: "hello"}, 0, 0, 10, 8)

	out := renderer.renderWindow(win)

	if !strings.Contains(out, "hello") {
		t.Errorf("rendered window %q missing content", out)
	}
}

// TestRenderWindowClipsTallContent guards against lipgloss Height/Width
// corrupting multi-line ANSI content (it used to interleave blank rows).
func TestRenderWindowClipsTallContent(t *testing.T) {
	t.Parallel()

	renderer := renderer{styles: config.Styles{}}

	lines := []string{"r1", "r2", "r3", "r4", "r5", "r6"}
	for i := range lines {
		lines[i] = "\x1b[31m" + lines[i] + "\x1b[0m"
	}

	win := NewWindow(1, staticModel{view: strings.Join(lines, "\n")}, 0, 0, 20, 6)

	out := renderer.renderWindow(win)

	got := strings.Split(out, "\n")

	// Window is 6 rows tall: title bar + 4 content rows + bottom border.
	if len(got) != 6 {
		t.Fatalf("rendered %d rows, want 6", len(got))
	}

	for i := 1; i <= 4; i++ {
		if strings.TrimSpace(stripANSI(got[i])) == "" {
			t.Errorf("content row %d is blank: %q", i, got[i])
		}
	}
}

func stripANSI(str string) string {
	var builder strings.Builder

	inEscape := false

	for i := range str {
		char := str[i]

		if inEscape {
			if char == 'm' {
				inEscape = false
			}

			continue
		}

		if char == 0x1b {
			inEscape = true

			continue
		}

		builder.WriteByte(char)
	}

	return builder.String()
}

func TestCanvasStampLeftClip(t *testing.T) {
	t.Parallel()

	cvs := newCanvas(10, 2)

	cvs.stamp(-2, 1, "hello")

	if !strings.Contains(cvs.lines[1], "llo") {
		t.Errorf("row 1 = %q, want it to contain the clipped 'llo'", cvs.lines[1])
	}

	if cvs.lines[0] != strings.Repeat(" ", 10) {
		t.Error("row 0 should remain blank")
	}
}

func TestCanvasStampRightClip(t *testing.T) {
	t.Parallel()

	cvs := newCanvas(5, 1)

	cvs.stamp(3, 0, "hello world")

	if cvs.lines[0] != "   he" {
		t.Errorf("row 0 = %q, want %q", cvs.lines[0], "   he")
	}
}
