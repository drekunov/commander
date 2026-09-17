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
