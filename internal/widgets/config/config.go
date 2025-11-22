package config

import (
	"github.com/charmbracelet/lipgloss"
)

type Config struct {
	ButtonStyle       lipgloss.Style
	ActiveButtonStyle lipgloss.Style
	DialogBoxStyle    lipgloss.Style
	TextStyle         lipgloss.Style
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
		Padding(0, 0).
		MarginTop(1)

	activeButtonStyle := buttonStyle.
		Foreground(lipgloss.Color("#FFF7DB")).
		Background(lipgloss.Color("#F25D94")).
		MarginRight(0).
		Underline(true)

	textStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF7DB")).
		Padding(0, 0).
		MarginTop(0).
		MarginBottom(0)

	Values = Config{
		ButtonStyle:       buttonStyle,
		ActiveButtonStyle: activeButtonStyle,
		DialogBoxStyle:    dialogBoxStyle,
		TextStyle:         textStyle,
	}
}
