package config

import (
	"embed"
	"os"
	"path"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

//go:embed styles/*.css
var defaultStyles embed.FS

const (
	stylesDir = "styles"
	cssSuffix = ".css"
)

// themeCSS assembles the embedded default theme from the per-type stylesheet
// files, concatenated in filename order. Every selector is defined once, so
// their order does not change the resulting theme.
func themeCSS() string {
	entries, err := defaultStyles.ReadDir(stylesDir)
	if err != nil {
		return ""
	}

	var builder strings.Builder

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), cssSuffix) {
			continue
		}

		data, readErr := defaultStyles.ReadFile(path.Join(stylesDir, entry.Name()))
		if readErr != nil {
			continue
		}

		builder.Write(data)
		builder.WriteString("\n")
	}

	return builder.String()
}

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
	FileDirectoryStyle        lipgloss.Style
	FileExecutableStyle       lipgloss.Style
	FileSymlinkStyle          lipgloss.Style
	FileImageStyle            lipgloss.Style
	FileArchiveStyle          lipgloss.Style
	FileSourceStyle           lipgloss.Style
	FileConfigStyle           lipgloss.Style
	MenuLabelStyle            lipgloss.Style
	MenuNumberStyle           lipgloss.Style
	MenuBackgroundStyle       lipgloss.Style
	TopBarStyle               lipgloss.Style
	MenuCaptionStyle          lipgloss.Style
	MenuCaptionActiveStyle    lipgloss.Style
	MenuHotkeyStyle           lipgloss.Style
	PulldownStyle             lipgloss.Style
	PulldownCursorStyle       lipgloss.Style
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
	cssData := themeCSS()

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
		FileDirectoryStyle:        ApplyStyle(stylesheet[".file-directory"]),
		FileExecutableStyle:       ApplyStyle(stylesheet[".file-executable"]),
		FileSymlinkStyle:          ApplyStyle(stylesheet[".file-symlink"]),
		FileImageStyle:            ApplyStyle(stylesheet[".file-image"]),
		FileArchiveStyle:          ApplyStyle(stylesheet[".file-archive"]),
		FileSourceStyle:           ApplyStyle(stylesheet[".file-source"]),
		FileConfigStyle:           ApplyStyle(stylesheet[".file-config"]),
		MenuLabelStyle:            ApplyStyle(stylesheet[".menu-label"]),
		MenuNumberStyle:           ApplyStyle(stylesheet[".menu-number"]),
		MenuBackgroundStyle:       ApplyStyle(stylesheet[".menu-background"]),
		TopBarStyle:               ApplyStyle(stylesheet[".menu-top-bar"]),
		MenuCaptionStyle:          ApplyStyle(stylesheet[".menu-caption"]),
		MenuCaptionActiveStyle:    ApplyStyle(stylesheet[".menu-caption-active"]),
		MenuHotkeyStyle:           ApplyStyle(stylesheet[".menu-hotkey"]),
		PulldownStyle:             ApplyStyle(stylesheet[".menu-pulldown"]),
		PulldownCursorStyle:       ApplyStyle(stylesheet[".menu-pulldown-cursor"]),
		WindowTitleStyle:          ApplyStyle(stylesheet[".window-title"]),
		WindowTitleUnfocusedStyle: ApplyStyle(stylesheet[".window-title-unfocused"]),
		WindowGripStyle:           ApplyStyle(stylesheet[".window-grip"]),
	}, nil
}
