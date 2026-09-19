//nolint:testpackage // tests exercise unexported menu state and helpers
package topmenu

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
)

func press(x, y int) tea.MouseMsg {
	return tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, X: x, Y: y}
}

func TestCaptionsRenderInOrder(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	row := ansi.Strip(model.View(100))

	last := -1

	for _, caption := range []string{"Left", "Files", "Commands", "Options", "Right"} {
		idx := strings.Index(row, caption)
		if idx < 0 {
			t.Fatalf("row %q missing caption %q", row, caption)
		}

		if idx < last {
			t.Errorf("caption %q appears out of order in %q", caption, row)
		}

		last = idx
	}

	if got := ansi.StringWidth(model.View(100)); got != 100 {
		t.Errorf("top row width = %d, want 100", got)
	}
}

func TestActiveCaptionUsesActiveStyle(t *testing.T) {
	t.Parallel()

	caption := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#111111")).
		Background(lipgloss.Color("#010101"))
	active := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#222222")).
		Background(lipgloss.Color("#020202"))
	hotkey := lipgloss.NewStyle().Foreground(lipgloss.Color("#333333"))

	model := New(config.Styles{
		MenuCaptionStyle:       caption,
		MenuCaptionActiveStyle: active,
		MenuHotkeyStyle:        hotkey,
	})
	model.Activate()
	model.menuIdx = 1

	if got := model.captionStyle(true).GetForeground(); got != active.GetForeground() {
		t.Errorf("active caption foreground = %v, want %v", got, active.GetForeground())
	}

	if got := model.captionStyle(false).GetForeground(); got != caption.GetForeground() {
		t.Errorf("inactive caption foreground = %v, want %v", got, caption.GetForeground())
	}

	if got := model.hotkeyStyle(true).GetForeground(); got != hotkey.GetForeground() {
		t.Errorf("hotkey foreground = %v, want %v", got, hotkey.GetForeground())
	}

	// The active highlight must span the whole caption, so the hotkey letter
	// keeps the active caption background instead of its own.
	if got := model.hotkeyStyle(true).GetBackground(); got != active.GetBackground() {
		t.Errorf("active hotkey background = %v, want active caption %v", got, active.GetBackground())
	}

	if got := model.hotkeyStyle(false).GetBackground(); got != caption.GetBackground() {
		t.Errorf("inactive hotkey background = %v, want caption %v", got, caption.GetBackground())
	}
}

func TestCaptionPaddingIsThemed(t *testing.T) {
	t.Parallel()

	caption := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#111111")).
		Background(lipgloss.Color("#010101")).
		Padding(0, 1, 0, 1)
	active := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#222222")).
		Background(lipgloss.Color("#020202")).
		Padding(0, 1, 0, 1)

	model := New(config.Styles{MenuCaptionStyle: caption, MenuCaptionActiveStyle: active})
	model.Activate()
	model.menuIdx = 0

	left, right := model.captionPadding()
	if left != 1 || right != 1 {
		t.Fatalf("caption padding = %d/%d, want 1/1", left, right)
	}

	wantWidth := len([]rune(menus[0].caption)) + 2
	if got := model.captionWidth(menus[0]); got != wantWidth {
		t.Errorf("caption width = %d, want %d", got, wantWidth)
	}

	if got := model.captionStart(1); got != wantWidth {
		t.Errorf("second caption start = %d, want %d", got, wantWidth)
	}

	// The active highlight spans the caption text plus its padding.
	seg := ansi.Strip(model.captionSegment(0, menus[0]))
	if got := ansi.StringWidth(seg); got != wantWidth {
		t.Errorf("caption segment width = %d, want %d (text plus padding)", got, wantWidth)
	}
}

func TestFallbackStylesRender(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 0
	model.openMenu()

	if model.View(80) == "" {
		t.Error("top row rendered empty with fallback styles")
	}

	if block, _ := model.PulldownBlock(80); block == "" {
		t.Error("pull-down rendered empty with fallback styles")
	}
}

func TestF9TogglesActive(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})

	if model.Active() {
		t.Fatal("new top menu should be inactive")
	}

	if _, handled := model.Update(tea.KeyMsg{Type: tea.KeyF9}); !handled {
		t.Error("F9 was not handled")
	}

	if !model.Active() {
		t.Error("F9 did not activate the top menu")
	}

	if _, handled := model.Update(tea.KeyMsg{Type: tea.KeyF9}); !handled {
		t.Error("second F9 was not handled")
	}

	if model.Active() {
		t.Error("second F9 did not deactivate the top menu")
	}
}

func TestArrowNavigationMovesMenus(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()

	model.Update(tea.KeyMsg{Type: tea.KeyRight})

	if model.menuIdx != 1 {
		t.Fatalf("menuIdx = %d after Right, want 1", model.menuIdx)
	}

	model.Update(tea.KeyMsg{Type: tea.KeyRight})

	if model.menuIdx != 2 {
		t.Fatalf("menuIdx = %d after second Right, want 2", model.menuIdx)
	}

	model.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if model.menuIdx != 1 {
		t.Fatalf("menuIdx = %d after Left, want 1", model.menuIdx)
	}

	model.menuIdx = 0
	model.Update(tea.KeyMsg{Type: tea.KeyLeft})

	if model.menuIdx != len(menus)-1 {
		t.Errorf("menuIdx = %d after Left at the first menu, want %d", model.menuIdx, len(menus)-1)
	}
}

func TestPulldownFollowsMenuChange(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 0
	model.openMenu()

	model.Update(tea.KeyMsg{Type: tea.KeyRight})

	if model.menuIdx != 1 || !model.Open() {
		t.Errorf("pull-down did not follow the menu change: menuIdx = %d, open = %v", model.menuIdx, model.Open())
	}
}

func TestDownOpensAndMovesInPulldown(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()

	model.Update(tea.KeyMsg{Type: tea.KeyDown})

	if !model.Open() || model.itemIdx != 0 {
		t.Fatalf("Down did not open the pull-down on its first item: open = %v, itemIdx = %d",
			model.Open(), model.itemIdx)
	}

	model.Update(tea.KeyMsg{Type: tea.KeyDown})

	if model.itemIdx != 1 {
		t.Errorf("itemIdx = %d after second Down, want 1", model.itemIdx)
	}

	model.Update(tea.KeyMsg{Type: tea.KeyUp})

	if model.itemIdx != 0 {
		t.Errorf("itemIdx = %d after Up, want 0", model.itemIdx)
	}

	model.Update(tea.KeyMsg{Type: tea.KeyUp})

	if model.itemIdx != 0 {
		t.Errorf("itemIdx = %d after Up at the top, want 0", model.itemIdx)
	}
}

func TestEnterOpensThenActivates(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 1 // Files, first item View

	cmd, _ := model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("Enter with a closed pull-down produced a command")
	}

	if !model.Open() {
		t.Fatal("Enter did not open the pull-down")
	}

	cmd, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("Enter on an item produced no command")
	}

	if model.Active() || model.Open() {
		t.Error("activating an item left the menu active or open")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != buttonbar.ActionView {
		t.Errorf("activated action = %v, want ActionView", msg.Action)
	}
}

func TestEscapeClosesThenDeactivates(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()

	model.Update(tea.KeyMsg{Type: tea.KeyDown})

	model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if model.Open() {
		t.Error("Escape left the pull-down open")
	}

	if !model.Active() {
		t.Error("Escape closed the pull-down but did not keep the bar active")
	}

	model.Update(tea.KeyMsg{Type: tea.KeyEsc})

	if model.Active() {
		t.Error("second Escape did not deactivate the bar")
	}
}

func TestHotkeyOpensMenuCaseInsensitive(t *testing.T) {
	t.Parallel()

	for _, typed := range []rune{'c', 'C'} {
		model := New(config.Styles{})
		model.Activate()

		model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{typed}})

		if model.menuIdx != 2 || !model.Open() {
			t.Errorf("hotkey %q did not open Commands: menuIdx = %d, open = %v",
				string(typed), model.menuIdx, model.Open())
		}
	}
}

func TestHotkeyIgnoredWhenInactive(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})

	if _, handled := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}}); handled {
		t.Error("inactive bar handled a hotkey")
	}

	if model.Active() || model.Open() {
		t.Error("inactive bar reacted to a hotkey")
	}
}

func TestFilesItemYieldsAction(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 1 // Files
	model.openMenu()
	model.itemIdx = 0 // View

	cmd := model.activateFocused()
	if cmd == nil {
		t.Fatal("activating View produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != buttonbar.ActionView {
		t.Errorf("View yielded action %v, want ActionView", msg.Action)
	}
}

func TestPlaceholderItemYieldsLabel(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 2 // Commands
	model.openMenu()
	model.itemIdx = 0 // Find file

	cmd := model.activateFocused()
	if cmd == nil {
		t.Fatal("activating a placeholder produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Action != 0 {
		t.Errorf("placeholder action = %v, want 0", msg.Action)
	}

	if msg.Label != labelFindFile {
		t.Errorf("placeholder label = %q, want Find file", msg.Label)
	}
}

func TestMouseClickMatchesKeyboard(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})

	captionX := model.captionStart(2) + 1

	// Click the Commands caption to open its pull-down.
	_, _ = model.HandleMouse(press(captionX, 0), 100)

	if model.menuIdx != 2 || !model.Open() {
		t.Fatalf("caption click did not open Commands: menuIdx = %d, open = %v", model.menuIdx, model.Open())
	}

	// The first item row is screen row 2 inside the box.
	cmd, consumed := model.HandleMouse(press(captionX, 2), 100)
	if !consumed {
		t.Fatal("item click was not consumed")
	}

	if cmd == nil {
		t.Fatal("item click produced no command")
	}

	msg, ok := cmd().(ActivateMsg)
	if !ok {
		t.Fatalf("cmd resolved to %T, want ActivateMsg", cmd())
	}

	if msg.Label != labelFindFile {
		t.Errorf("item click activated %q, want Find file", msg.Label)
	}
}

func TestMouseClickOutsideBoxNotConsumed(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 0 // Left, its box starts at column 0
	model.openMenu()

	if _, consumed := model.HandleMouse(press(40, 2), 80); consumed {
		t.Error("a click to the right of the pull-down box was consumed")
	}
}

func TestPulldownBlockIsCompactBox(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 1 // Files
	model.openMenu()

	block, startX := model.PulldownBlock(80)
	lines := strings.Split(block, "\n")

	if len(lines) != len(menus[1].items)+2 {
		t.Fatalf("pull-down line count = %d, want %d", len(lines), len(menus[1].items)+2)
	}

	if startX != model.captionStart(1) {
		t.Errorf("pull-down start column = %d, want %d", startX, model.captionStart(1))
	}

	boxWidth := ansi.StringWidth(lines[0])
	if boxWidth >= 80 {
		t.Errorf("pull-down box width = %d, want a compact box narrower than the terminal", boxWidth)
	}

	for idx, line := range lines {
		if got := ansi.StringWidth(line); got != boxWidth {
			t.Errorf("line %d width = %d, want box width %d", idx, got, boxWidth)
		}
	}

	if !strings.Contains(ansi.Strip(lines[0]), "┌") {
		t.Error("pull-down has no top border")
	}

	if !strings.Contains(ansi.Strip(lines[1]), "View") {
		t.Error("first pull-down row does not show View")
	}
}

func TestPulldownCursorExcludesFrame(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 1 // Files
	model.openMenu()

	// The item text carries no frame characters, so the cursor style applied
	// to it can never select the frame.
	text := model.pulldownItemText(10, "View")
	if strings.Contains(text, "│") {
		t.Errorf("item text %q contains frame characters", text)
	}

	if got := ansi.StringWidth(text); got != 6 {
		t.Errorf("item text width = %d, want 6 (box width minus frame)", got)
	}

	row := model.pulldownItemRow(10, "View", model.itemStyle(0))
	if got := ansi.Strip(row); got != "│ View   │" {
		t.Errorf("row = %q, want %q", got, "│ View   │")
	}
}

func TestPulldownFocusedItemUsesCursorStyle(t *testing.T) {
	t.Parallel()

	pulldown := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	cursor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF00FF"))

	model := New(config.Styles{PulldownStyle: pulldown, PulldownCursorStyle: cursor})
	model.Activate()
	model.menuIdx = 1
	model.openMenu()

	if got := model.itemStyle(0).GetForeground(); got != cursor.GetForeground() {
		t.Errorf("focused item foreground = %v, want cursor %v", got, cursor.GetForeground())
	}

	if got := model.itemStyle(1).GetForeground(); got != pulldown.GetForeground() {
		t.Errorf("unfocused item foreground = %v, want pull-down %v", got, pulldown.GetForeground())
	}
}

func TestPulldownTruncatesOnNarrowTerminal(t *testing.T) {
	t.Parallel()

	model := New(config.Styles{})
	model.Activate()
	model.menuIdx = 2 // Commands has the long "Compare directories" label
	model.openMenu()

	block, _ := model.PulldownBlock(20)

	for idx, line := range strings.Split(block, "\n") {
		if got := ansi.StringWidth(line); got > 20 {
			t.Errorf("line %d width = %d, want at most 20", idx, got)
		}
	}
}
