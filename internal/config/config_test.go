package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesNoOverride(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.DialogBoxStyle.GetBorderStyle() != lipgloss.DoubleBorder() {
		t.Error("expected embedded default double border when no override present")
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesOverride(t *testing.T) {
	dir := t.TempDir()

	override := []byte(".dialog-box { color: #FF0000; }")

	err := os.WriteFile(filepath.Join(dir, "styles.css"), override, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.DialogBoxStyle.GetForeground() != lipgloss.Color("#FF0000") {
		t.Error("expected override color to be applied")
	}
}

func TestParseCSS(t *testing.T) {
	t.Parallel()

	sheet := config.ParseCSS(".button { color: #000000; padding: 0 1 0 1; }")

	props, ok := sheet[".button"]
	if !ok {
		t.Fatal("expected .button selector")
	}

	if props["color"] != "#000000" {
		t.Errorf("color = %q, want #000000", props["color"])
	}
}

func TestApplyStyle(t *testing.T) {
	t.Parallel()

	style := config.ApplyStyle(map[string]string{
		"color":            "#123456",
		"background-color": "#654321",
		"padding":          "0 1 0 1",
		"border-style":     "rounded",
	})

	if style.GetForeground() != lipgloss.Color("#123456") {
		t.Error("foreground color not applied")
	}

	if style.GetBackground() != lipgloss.Color("#654321") {
		t.Error("background color not applied")
	}

	if style.GetPaddingLeft() != 1 || style.GetPaddingRight() != 1 {
		t.Error("horizontal padding not applied")
	}

	if style.GetBorderStyle() != lipgloss.RoundedBorder() {
		t.Error("border style not applied")
	}
}

func TestApplyStyleBold(t *testing.T) {
	t.Parallel()

	bold := config.ApplyStyle(map[string]string{"font-weight": "bold"})
	if !bold.GetBold() {
		t.Error("font-weight: bold did not apply bold")
	}

	plain := config.ApplyStyle(map[string]string{"color": "#000000"})
	if plain.GetBold() {
		t.Error("style without font-weight should not be bold")
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesCursorAndMenuLabel(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name  string
		style lipgloss.Style
	}{
		{name: "cursor", style: styles.CursorStyle},
		{name: "menu-label", style: styles.MenuLabelStyle},
	}

	for _, testCase := range cases {
		if testCase.style.GetForeground() != lipgloss.Color("#000000") {
			t.Errorf("%s foreground = %v, want #000000", testCase.name, testCase.style.GetForeground())
		}

		if testCase.style.GetBackground() != lipgloss.Color("#008080") {
			t.Errorf("%s background = %v, want #008080", testCase.name, testCase.style.GetBackground())
		}

		if !testCase.style.GetBold() {
			t.Errorf("%s style is not bold", testCase.name)
		}
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesMenuNumber(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.MenuNumberStyle.GetForeground() != lipgloss.Color("#C0C0C0") {
		t.Errorf("menu-number foreground = %v, want #C0C0C0", styles.MenuNumberStyle.GetForeground())
	}

	if styles.MenuNumberStyle.GetBackground() != lipgloss.Color("#000000") {
		t.Errorf("menu-number background = %v, want #000000", styles.MenuNumberStyle.GetBackground())
	}

	if !styles.MenuNumberStyle.GetBold() {
		t.Error("menu-number style is not bold")
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesMenuBackground(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.MenuBackgroundStyle.GetForeground() != lipgloss.Color("#C0C0C0") {
		t.Errorf("menu-background foreground = %v, want #C0C0C0", styles.MenuBackgroundStyle.GetForeground())
	}

	if styles.MenuBackgroundStyle.GetBackground() != lipgloss.Color("#000000") {
		t.Errorf("menu-background background = %v, want #000000", styles.MenuBackgroundStyle.GetBackground())
	}

	if !styles.MenuBackgroundStyle.GetBold() {
		t.Error("menu-background style is not bold")
	}
}

// chromeCase describes the expected attributes of one theme class.
type chromeCase struct {
	name    string
	style   lipgloss.Style
	fg      lipgloss.TerminalColor
	bg      lipgloss.TerminalColor
	bold    bool
	padding int
}

func assertChromeStyles(t *testing.T, cases []chromeCase) {
	t.Helper()

	for _, testCase := range cases {
		if testCase.fg != nil && testCase.style.GetForeground() != testCase.fg {
			t.Errorf("%s foreground = %v, want %v", testCase.name, testCase.style.GetForeground(), testCase.fg)
		}

		if testCase.bg != nil && testCase.style.GetBackground() != testCase.bg {
			t.Errorf("%s background = %v, want %v", testCase.name, testCase.style.GetBackground(), testCase.bg)
		}

		if testCase.style.GetBold() != testCase.bold {
			t.Errorf("%s bold = %v, want %v", testCase.name, testCase.style.GetBold(), testCase.bold)
		}

		if got := testCase.style.GetPaddingLeft(); got != testCase.padding {
			t.Errorf("%s padding left = %d, want %d", testCase.name, got, testCase.padding)
		}

		if got := testCase.style.GetPaddingRight(); got != testCase.padding {
			t.Errorf("%s padding right = %d, want %d", testCase.name, got, testCase.padding)
		}
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesDialogChrome(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	assertChromeStyles(t, []chromeCase{
		{
			name:    "dialog-cursor",
			style:   styles.DialogCursorStyle,
			fg:      lipgloss.Color("#000000"),
			bg:      lipgloss.Color("#00AAAA"),
			bold:    true,
			padding: 1,
		},
		{
			name:    "dialog-option",
			style:   styles.DialogOptionStyle,
			fg:      lipgloss.Color("#00AAAA"),
			bg:      lipgloss.Color("#000080"),
			padding: 1,
		},
	})
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesWindowChrome(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	assertChromeStyles(t, []chromeCase{
		{
			name:  "window-title",
			style: styles.WindowTitleStyle,
			fg:    lipgloss.Color("#FFFFFF"),
			bg:    lipgloss.Color("#000080"),
			bold:  true,
		},
		{
			name:  "window-title-unfocused",
			style: styles.WindowTitleUnfocusedStyle,
			fg:    lipgloss.Color("#00AAAA"),
			bg:    lipgloss.Color("#000080"),
		},
		{
			name:  "window-grip",
			style: styles.WindowGripStyle,
			fg:    lipgloss.Color("#00AAAA"),
			bg:    lipgloss.Color("#000080"),
		},
	})
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesDialogAndWindowOverride(t *testing.T) {
	dir := t.TempDir()

	override := []byte(".dialog-cursor { background-color: #123456; }\n.window-grip { color: #ABCDEF; }")

	err := os.WriteFile(filepath.Join(dir, "styles.css"), override, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.DialogCursorStyle.GetBackground() != lipgloss.Color("#123456") {
		t.Errorf("dialog-cursor background = %v, want #123456", styles.DialogCursorStyle.GetBackground())
	}

	if styles.WindowGripStyle.GetForeground() != lipgloss.Color("#ABCDEF") {
		t.Errorf("window-grip foreground = %v, want #ABCDEF", styles.WindowGripStyle.GetForeground())
	}
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesTopMenu(t *testing.T) {
	t.Chdir(t.TempDir())

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	assertChromeStyles(t, []chromeCase{
		{
			name:  "menu-top-bar",
			style: styles.TopBarStyle,
			fg:    lipgloss.Color("#C0C0C0"),
			bg:    lipgloss.Color("#000000"),
		},
		{
			name:    "menu-caption",
			style:   styles.MenuCaptionStyle,
			fg:      lipgloss.Color("#C0C0C0"),
			bg:      lipgloss.Color("#000000"),
			padding: 1,
		},
		{
			name:    "menu-caption-active",
			style:   styles.MenuCaptionActiveStyle,
			fg:      lipgloss.Color("#000000"),
			bg:      lipgloss.Color("#00AAAA"),
			padding: 1,
		},
		{
			name:  "menu-hotkey",
			style: styles.MenuHotkeyStyle,
			fg:    lipgloss.Color("#FFFF00"),
			bg:    lipgloss.Color("#000000"),
			bold:  true,
		},
		{
			name:  "menu-pulldown",
			style: styles.PulldownStyle,
			fg:    lipgloss.Color("#C0C0C0"),
			bg:    lipgloss.Color("#000080"),
		},
		{
			name:  "menu-pulldown-cursor",
			style: styles.PulldownCursorStyle,
			fg:    lipgloss.Color("#000000"),
			bg:    lipgloss.Color("#00AAAA"),
			bold:  true,
		},
	})
}

//nolint:paralleltest // t.Chdir mutates the process-global working directory
func TestLoadStylesTopMenuOverride(t *testing.T) {
	dir := t.TempDir()

	override := []byte(".menu-caption { color: #111111; }\n.menu-pulldown-cursor { background-color: #222222; }")

	err := os.WriteFile(filepath.Join(dir, "styles.css"), override, 0o600)
	if err != nil {
		t.Fatal(err)
	}

	t.Chdir(dir)

	styles, err := config.LoadStyles()
	if err != nil {
		t.Fatal(err)
	}

	if styles.MenuCaptionStyle.GetForeground() != lipgloss.Color("#111111") {
		t.Errorf("menu-caption foreground = %v, want #111111", styles.MenuCaptionStyle.GetForeground())
	}

	if styles.PulldownCursorStyle.GetBackground() != lipgloss.Color("#222222") {
		t.Errorf("menu-pulldown-cursor background = %v, want #222222", styles.PulldownCursorStyle.GetBackground())
	}
}
