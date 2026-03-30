package config

import (
	_ "embed"
	"os"

	"github.com/charmbracelet/lipgloss"
)

//go:embed styles.css
var defaultCSS string

type Config struct {
	ButtonStyle        lipgloss.Style
	ActiveButtonStyle  lipgloss.Style
	PressedButtonStyle lipgloss.Style
	DialogBoxStyle     lipgloss.Style
	TextStyle          lipgloss.Style
	InputStyle         lipgloss.Style
	TableStyle         lipgloss.Style
}

var Values Config

func init() {
	cssData := defaultCSS

	data, err := os.ReadFile("styles.css")
	if err == nil {
		cssData = string(data)
	}

	stylesheet := ParseCSS(cssData)

	Values = Config{
		ButtonStyle:        ApplyStyle(stylesheet[".button"]),
		ActiveButtonStyle:  ApplyStyle(stylesheet[".active-button"]),
		PressedButtonStyle: ApplyStyle(stylesheet[".pressed-button"]),
		DialogBoxStyle:     ApplyStyle(stylesheet[".dialog-box"]),
		TextStyle:          ApplyStyle(stylesheet[".text"]),
		InputStyle:         ApplyStyle(stylesheet[".input"]),
		TableStyle:         ApplyStyle(stylesheet[".table"]),
	}
}
