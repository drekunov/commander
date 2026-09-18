//nolint:testpackage // tests exercise unexported model state and helpers
package buttonbar

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	if updated.pressed != ActionCopy {
		t.Error("pressed should persist until the next Update, not be cleared by View")
	}

	updated, _ = updated.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	if updated.pressed != 0 {
		t.Error("pressed was not cleared at the start of the next Update")
	}
}

func TestPressedSegmentKeepsBarBackground(t *testing.T) {
	t.Parallel()

	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))
	pressed := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Background(lipgloss.Color("#00AAAA"))

	model := New(config.Styles{
		ButtonBarStyle:     bar,
		PressedButtonStyle: pressed,
	})
	model.pressed = ActionCopy

	got := model.styleFor(0, ActionCopy)

	if got.GetBackground() != bar.GetBackground() {
		t.Errorf("pressed segment background = %v, want bar background %v", got.GetBackground(), bar.GetBackground())
	}

	if got.GetForeground() != pressed.GetForeground() {
		t.Errorf("pressed segment foreground = %v, want pressed foreground %v", got.GetForeground(), pressed.GetForeground())
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

	if got := Action(0).Label(); got != "" {
		t.Errorf("Action(0).Label() = %q, want empty", got)
	}

	if got := Action(99).Name(); got != "" {
		t.Errorf("Action(99).Name() = %q, want empty", got)
	}
}

func TestDefaultLabelUsesMenuLabelOnStrip(t *testing.T) {
	t.Parallel()

	menuLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#C0C0C0"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{ButtonBarStyle: bar, MenuLabelStyle: menuLabel})

	got := model.styleFor(0, actions[0].action)

	if got.GetBackground() != menuLabel.GetBackground() {
		t.Errorf("default label background = %v, want menu-label %v", got.GetBackground(), menuLabel.GetBackground())
	}

	if got.GetForeground() != menuLabel.GetForeground() {
		t.Errorf("default label foreground = %v, want menu-label %v", got.GetForeground(), menuLabel.GetForeground())
	}
}

func TestFocusedLabelUsesActiveStyleOnStrip(t *testing.T) {
	t.Parallel()

	active := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#000080"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{ButtonBarStyle: bar, ActiveButtonStyle: active})
	model.SetFocused(true)

	got := model.styleFor(0, actions[0].action)

	if got.GetBackground() != bar.GetBackground() {
		t.Errorf("focused label background = %v, want strip %v", got.GetBackground(), bar.GetBackground())
	}

	if got.GetForeground() != active.GetForeground() {
		t.Errorf("focused label foreground = %v, want active %v", got.GetForeground(), active.GetForeground())
	}
}

func TestNumberAndNameUseSeparateStyles(t *testing.T) {
	t.Parallel()

	menuNumber := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0")).
		Background(lipgloss.Color("#000000"))
	menuLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#000000")).
		Background(lipgloss.Color("#C0C0C0"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{ButtonBarStyle: bar, MenuLabelStyle: menuLabel, MenuNumberStyle: menuNumber})

	number := model.numberStyle()
	if number.GetForeground() != menuNumber.GetForeground() {
		t.Errorf("number foreground = %v, want menu-number %v", number.GetForeground(), menuNumber.GetForeground())
	}

	if number.GetBackground() != menuNumber.GetBackground() {
		t.Errorf("number background = %v, want menu-number %v", number.GetBackground(), menuNumber.GetBackground())
	}

	name := model.styleFor(0, actions[0].action)
	if name.GetForeground() != menuLabel.GetForeground() {
		t.Errorf("name foreground = %v, want menu-label %v", name.GetForeground(), menuLabel.GetForeground())
	}

	if name.GetBackground() != menuLabel.GetBackground() {
		t.Errorf("name background = %v, want menu-label %v", name.GetBackground(), menuLabel.GetBackground())
	}
}

func TestNumberKeepsStyleOnFocusedAndPressed(t *testing.T) {
	t.Parallel()

	menuNumber := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0")).
		Background(lipgloss.Color("#000000"))
	active := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#000080"))
	pressed := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Background(lipgloss.Color("#000080"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{
		ButtonBarStyle:     bar,
		MenuNumberStyle:    menuNumber,
		ActiveButtonStyle:  active,
		PressedButtonStyle: pressed,
	})
	model.SetFocused(true)

	number := model.numberStyle()
	if number.GetForeground() != menuNumber.GetForeground() || number.GetBackground() != menuNumber.GetBackground() {
		t.Errorf("focused number style = %v/%v, want menu-number %v/%v",
			number.GetForeground(), number.GetBackground(), menuNumber.GetForeground(), menuNumber.GetBackground())
	}

	if name := model.styleFor(0, actions[0].action); name.GetForeground() != active.GetForeground() {
		t.Errorf("focused name foreground = %v, want active %v", name.GetForeground(), active.GetForeground())
	}

	model.pressed = actions[1].action

	number = model.numberStyle()
	if number.GetForeground() != menuNumber.GetForeground() || number.GetBackground() != menuNumber.GetBackground() {
		t.Errorf("pressed number style = %v/%v, want menu-number %v/%v",
			number.GetForeground(), number.GetBackground(), menuNumber.GetForeground(), menuNumber.GetBackground())
	}

	if name := model.styleFor(1, actions[1].action); name.GetForeground() != pressed.GetForeground() {
		t.Errorf("pressed name foreground = %v, want pressed %v", name.GetForeground(), pressed.GetForeground())
	}
}

func TestButtonsHaveEqualWidth(t *testing.T) {
	t.Parallel()

	for _, width := range []int{80, 90, 100, 121} {
		model := New(config.Styles{})
		model = resize(model, width)

		if got := model.buttonWidth(); got != naturalButtonWidth {
			t.Errorf("width %d: buttonWidth = %d, want %d", width, got, naturalButtonWidth)
		}

		row := model.View()
		if got := ansi.StringWidth(row); got != width {
			t.Errorf("width %d: row width = %d, want %d", width, got, width)
		}

		first := model.gapWidth(0)
		for i := 1; i < buttonCount-1; i++ {
			if diff := model.gapWidth(i) - first; diff < -1 || diff > 1 {
				t.Errorf("width %d: gap %d = %d differs from gap 0 = %d",
					width, i, model.gapWidth(i), first)
			}
		}
	}
}

func TestNarrowTerminalShrinksButtonsEqually(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model = resize(model, 60)

	if got := model.buttonWidth(); got != 6 {
		t.Fatalf("buttonWidth = %d, want 6", got)
	}

	row := model.View()
	if got := ansi.StringWidth(row); got != 60 {
		t.Errorf("row width = %d, want 60", got)
	}

	if strings.Contains(ansi.Strip(row), "6RenMov") {
		t.Error("expected the RenMov label to be truncated at width 60")
	}
}

func TestMouseMappingWithGaps(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model = resize(model, 100)

	buttonWidth := model.buttonWidth()

	start := 0
	for i := range 5 {
		start += buttonWidth + model.gapWidth(i)
	}

	_, cmd := model.Update(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      start,
		Y:      0,
	})
	if cmd == nil {
		t.Fatal("click at the start of button 5 produced no command")
	}

	if msg, ok := cmd().(ActivateMsg); !ok || msg.Action != actions[5].action {
		t.Errorf("click at column %d activated %v, want %v", start, cmd(), actions[5].action)
	}

	if gap := model.gapWidth(5); gap > 0 {
		_, cmd = model.Update(tea.MouseMsg{
			Action: tea.MouseActionPress,
			Button: tea.MouseButtonLeft,
			X:      start + buttonWidth,
			Y:      0,
		})

		if msg, ok := cmd().(ActivateMsg); !ok || msg.Action != actions[5].action {
			t.Errorf("click in the gap after button 5 activated %v, want %v", cmd(), actions[5].action)
		}
	}
}

func TestBarBackgroundUsesMenuBackgroundStyle(t *testing.T) {
	t.Parallel()

	menuBackground := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0")).
		Background(lipgloss.Color("#000000"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{ButtonBarStyle: bar, MenuBackgroundStyle: menuBackground})

	base := model.barBackground()
	if base.GetBackground() != menuBackground.GetBackground() {
		t.Errorf("bar background = %v, want menu-background %v", base.GetBackground(), menuBackground.GetBackground())
	}

	if base.GetForeground() != menuBackground.GetForeground() {
		t.Errorf("bar background foreground = %v, want menu-background %v",
			base.GetForeground(), menuBackground.GetForeground())
	}

	// An unset menu-background falls back to the button-bar strip.
	fallback := New(config.Styles{ButtonBarStyle: bar}).barBackground()
	if fallback.GetBackground() != bar.GetBackground() {
		t.Errorf("bar background fallback = %v, want button-bar %v", fallback.GetBackground(), bar.GetBackground())
	}
}

func TestNumberAndNameInheritBarBackground(t *testing.T) {
	t.Parallel()

	menuBackground := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#C0C0C0")).
		Background(lipgloss.Color("#000000"))

	// With no number or label style, both inherit the menu-background base.
	model := New(config.Styles{MenuBackgroundStyle: menuBackground})

	if got := model.numberStyle(); got.GetBackground() != menuBackground.GetBackground() {
		t.Errorf("number base background = %v, want menu-background %v", got.GetBackground(), menuBackground.GetBackground())
	}

	if got := model.styleFor(0, actions[0].action); got.GetBackground() != menuBackground.GetBackground() {
		t.Errorf("name base background = %v, want menu-background %v", got.GetBackground(), menuBackground.GetBackground())
	}
}

func TestFocusedAndPressedKeepMenuBackgroundBase(t *testing.T) {
	t.Parallel()

	menuBackground := lipgloss.NewStyle().Background(lipgloss.Color("#000000"))
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	pressed := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	bar := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAAA")).
		Background(lipgloss.Color("#000080"))

	model := New(config.Styles{
		ButtonBarStyle:      bar,
		MenuBackgroundStyle: menuBackground,
		ActiveButtonStyle:   active,
		PressedButtonStyle:  pressed,
	})
	model.SetFocused(true)

	focused := model.styleFor(0, actions[0].action)
	if focused.GetBackground() != menuBackground.GetBackground() {
		t.Errorf("focused background = %v, want menu-background %v", focused.GetBackground(), menuBackground.GetBackground())
	}

	if focused.GetForeground() != active.GetForeground() {
		t.Errorf("focused foreground = %v, want active %v", focused.GetForeground(), active.GetForeground())
	}

	model.pressed = actions[1].action

	pressedStyle := model.styleFor(1, actions[1].action)
	if pressedStyle.GetBackground() != menuBackground.GetBackground() {
		t.Errorf("pressed background = %v, want menu-background %v",
			pressedStyle.GetBackground(), menuBackground.GetBackground())
	}

	if pressedStyle.GetForeground() != pressed.GetForeground() {
		t.Errorf("pressed foreground = %v, want pressed %v", pressedStyle.GetForeground(), pressed.GetForeground())
	}
}

func TestClearPressedRestoresNormalStyle(t *testing.T) {
	t.Parallel()

	menuLabel := lipgloss.NewStyle().Foreground(lipgloss.Color("#00AAAA"))
	pressed := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))

	model := New(config.Styles{MenuLabelStyle: menuLabel, PressedButtonStyle: pressed})
	model.pressed = actions[1].action

	if got := model.styleFor(1, actions[1].action); got.GetForeground() != pressed.GetForeground() {
		t.Fatalf("pressed foreground = %v, want pressed %v", got.GetForeground(), pressed.GetForeground())
	}

	model.ClearPressed()

	if got := model.styleFor(1, actions[1].action); got.GetForeground() != menuLabel.GetForeground() {
		t.Errorf("after ClearPressed foreground = %v, want menu-label %v", got.GetForeground(), menuLabel.GetForeground())
	}
}
