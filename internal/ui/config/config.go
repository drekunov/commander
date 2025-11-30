package config

import (
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	ButtonStyle        lipgloss.Style
	ActiveButtonStyle  lipgloss.Style
	PressedButtonStyle lipgloss.Style
	DialogBoxStyle     lipgloss.Style
	TextStyle          lipgloss.Style
	InputStyle         lipgloss.Style
}

var Values Config

func init() {
	dialogBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#874BFD")).
		Padding(0, 1, 0, 1).
		BorderTop(true).
		BorderLeft(true).
		BorderRight(true).
		BorderBottom(true)

	buttonStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#888B7E")).
		Padding(0, 1, 0, 1).
		Margin(1, 1, 1, 1)

	activeButtonStyle := buttonStyle.
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		Underline(true)

	pressedButtonStyle := activeButtonStyle.
		Background(lipgloss.Color("#725D94"))

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Padding(0, 0).
		Margin(1, 1, 0, 1)

	inputStyle := textStyle.
		Foreground(lipgloss.Color("#874BFD")).
		Padding(0, 0, 0, 0).
		Margin(1, 1, 0, 1)

	Values = Config{
		ButtonStyle:        buttonStyle,
		ActiveButtonStyle:  activeButtonStyle,
		PressedButtonStyle: pressedButtonStyle,
		DialogBoxStyle:     dialogBoxStyle,
		TextStyle:          textStyle,
		InputStyle:         inputStyle,
	}
}
