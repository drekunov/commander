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
	css := defaultCSS
	if data, err := os.ReadFile("styles.css"); err == nil {
		css = string(data)
	}

	ss := ParseCSS(css)

	Values = Config{
		ButtonStyle:        ApplyStyle(ss[".button"]),
		ActiveButtonStyle:  ApplyStyle(ss[".active-button"]),
		PressedButtonStyle: ApplyStyle(ss[".pressed-button"]),
		DialogBoxStyle:     ApplyStyle(ss[".dialog-box"]),
		TextStyle:          ApplyStyle(ss[".text"]),
		InputStyle:         ApplyStyle(ss[".input"]),
		TableStyle:         ApplyStyle(ss[".table"]),
	}
}