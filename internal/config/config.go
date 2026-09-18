package config

import (
	_ "embed"
	"os"

	"github.com/charmbracelet/lipgloss"
)

//go:embed styles.css
var defaultCSS string

type Styles struct {
	ButtonStyle               lipgloss.Style
	ActiveButtonStyle         lipgloss.Style
	PressedButtonStyle        lipgloss.Style
	DialogBoxStyle            lipgloss.Style
	DialogCursorStyle         lipgloss.Style
	DialogOptionStyle         lipgloss.Style
	TextStyle                 lipgloss.Style
	InputStyle                lipgloss.Style
	TableStyle                lipgloss.Style
	ButtonBarStyle            lipgloss.Style
	CursorStyle               lipgloss.Style
	MenuLabelStyle            lipgloss.Style
	MenuNumberStyle           lipgloss.Style
	MenuBackgroundStyle       lipgloss.Style
	WindowTitleStyle          lipgloss.Style
	WindowTitleUnfocusedStyle lipgloss.Style
	WindowGripStyle           lipgloss.Style
}

// DialogBase returns a style carrying the dialog frame's background. Dialog
// body elements inherit it so their content sits on the frame instead of the
// terminal's default background.
func (s Styles) DialogBase() lipgloss.Style {
	return lipgloss.NewStyle().Background(s.DialogBoxStyle.GetBackground())
}

// LoadStyles builds the theme, applying a styles.css override from the
// process working directory when present and falling back to the embedded
// default otherwise.
func LoadStyles() (Styles, error) {
	cssData := defaultCSS

	data, err := os.ReadFile("styles.css")
	if err == nil {
		cssData = string(data)
	}

	stylesheet := ParseCSS(cssData)

	return Styles{
		ButtonStyle:               ApplyStyle(stylesheet[".button"]),
		ActiveButtonStyle:         ApplyStyle(stylesheet[".active-button"]),
		PressedButtonStyle:        ApplyStyle(stylesheet[".pressed-button"]),
		DialogBoxStyle:            ApplyStyle(stylesheet[".dialog-box"]),
		DialogCursorStyle:         ApplyStyle(stylesheet[".dialog-cursor"]),
		DialogOptionStyle:         ApplyStyle(stylesheet[".dialog-option"]),
		TextStyle:                 ApplyStyle(stylesheet[".text"]),
		InputStyle:                ApplyStyle(stylesheet[".input"]),
		TableStyle:                ApplyStyle(stylesheet[".table"]),
		ButtonBarStyle:            ApplyStyle(stylesheet[".button-bar"]),
		CursorStyle:               ApplyStyle(stylesheet[".cursor"]),
		MenuLabelStyle:            ApplyStyle(stylesheet[".menu-label"]),
		MenuNumberStyle:           ApplyStyle(stylesheet[".menu-number"]),
		MenuBackgroundStyle:       ApplyStyle(stylesheet[".menu-background"]),
		WindowTitleStyle:          ApplyStyle(stylesheet[".window-title"]),
		WindowTitleUnfocusedStyle: ApplyStyle(stylesheet[".window-title-unfocused"]),
		WindowGripStyle:           ApplyStyle(stylesheet[".window-grip"]),
	}, nil
}
