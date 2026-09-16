//nolint:testpackage // tests exercise unexported model state and helpers
package buttonbar

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
)

func resize(m *Model, width int) *Model {
	updated, _ := m.Update(tea.WindowSizeMsg{Width: width, Height: 30})

	return updated
}

func TestRenderTenButtonsInOrder(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model = resize(model, 100)

	row := ansi.Strip(model.View())

	for _, def := range actions {
		if !strings.Contains(row, def.label) {
			t.Errorf("row %q missing button %q", row, def.label)
		}
	}
}

func TestRowFillsExactWidth(t *testing.T) {
	t.Parallel()

	for _, width := range []int{30, 41, 60, 80, 100, 121} {
		model := New(config.Styles{})
		model = resize(model, width)

		row := model.View()

		if got := ansi.StringWidth(row); got != width {
			t.Errorf("width %d: row width = %d, want %d", width, got, width)
		}

		if strings.Count(row, "\n") != 0 {
			t.Errorf("width %d: row wraps to multiple lines", width)
		}
	}
}

func TestActionForKey(t *testing.T) {
	t.Parallel()

	for i, def := range actions {
		key := tea.KeyF1 - tea.KeyType(i)

		got, ok := ActionForKey(key)
		if !ok {
			t.Fatalf("ActionForKey(%v) not found", key)
		}

		if got != def.action {
			t.Errorf("ActionForKey(%v) = %v, want %v", key, got, def.action)
		}
	}

	if _, ok := ActionForKey(tea.KeyF12); ok {
		t.Error("ActionForKey(F12) should not match a bar button")
	}
}

func TestFunctionKeyActivation(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model = resize(model, 100)

	updated, cmd := model.Update(tea.KeyMsg{Type: tea.KeyF5})
	if cmd == nil {
		t.Fatal("F5 produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != ActionCopy {
		t.Errorf("F5 activated %v, want ActionCopy", msg.Action)
	}

	if updated.pressed != ActionCopy {
		t.Errorf("pressed = %v, want ActionCopy after activation", updated.pressed)
	}

	updated.View()

	if updated.pressed != 0 {
		t.Error("pressed was not cleared after rendering the flash frame")
	}
}

func TestKeyboardNavigationWhenFocused(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.SetFocused(true)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})

	if model.focusedIdx != 2 {
		t.Fatalf("focusedIdx = %d, want 2 after two Rights", model.focusedIdx)
	}

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on a focused bar produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != actions[2].action {
		t.Errorf("Enter activated %v, want %v", msg.Action, actions[2].action)
	}
}

func TestLeftWrapsToLastButton(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.SetFocused(true)

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if model.focusedIdx != len(actions)-1 {
		t.Errorf("focusedIdx = %d, want last button %d", model.focusedIdx, len(actions)-1)
	}
}

func TestNavigationIgnoredWhenNotFocused(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})

	_, cmd := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("Enter produced a command while the bar is not focused")
	}

	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRight})

	if model.focusedIdx != 0 {
		t.Errorf("focusedIdx changed to %d while not focused", model.focusedIdx)
	}
}

func TestMouseClickActivatesButton(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model = resize(model, 100)

	// 5Copy starts at column 40 in a 100-wide row (ten 10-wide segments).
	_, cmd := model.Update(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      44,
		Y:      0,
	})

	if cmd == nil {
		t.Fatal("mouse press produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != ActionCopy {
		t.Errorf("click at column 44 activated %v, want ActionCopy", msg.Action)
	}
}

func TestLabelAndName(t *testing.T) {
	t.Parallel()

	if got := ActionCopy.Label(); got != "5Copy" {
		t.Errorf("ActionCopy.Label() = %q, want %q", got, "5Copy")
	}

	if got := ActionCopy.Name(); got != "Copy" {
		t.Errorf("ActionCopy.Name() = %q, want %q", got, "Copy")
	}
}
