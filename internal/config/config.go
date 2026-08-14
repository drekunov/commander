package config

import (
	_ "embed"
	"os"

	"github.com/charmbracelet/lipgloss"
)

//go:embed styles.css
var defaultCSS string

type Styles struct {
	ButtonStyle        lipgloss.Style
	ActiveButtonStyle  lipgloss.Style
	PressedButtonStyle lipgloss.Style
	DialogBoxStyle     lipgloss.Style
	TextStyle          lipgloss.Style
	InputStyle         lipgloss.Style
	TableStyle         lipgloss.Style
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
		ButtonStyle:        ApplyStyle(stylesheet[".button"]),
		ActiveButtonStyle:  ApplyStyle(stylesheet[".active-button"]),
		PressedButtonStyle: ApplyStyle(stylesheet[".pressed-button"]),
		DialogBoxStyle:     ApplyStyle(stylesheet[".dialog-box"]),
		TextStyle:          ApplyStyle(stylesheet[".text"]),
		InputStyle:         ApplyStyle(stylesheet[".input"]),
		TableStyle:         ApplyStyle(stylesheet[".table"]),
	}, nil
}
