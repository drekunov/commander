package dialogs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

// dialogBase provides the shared footer and visibility shell used by every
// dialog widget in this package. The window manager draws the single frame and
// title around the dialog's body.
type dialogBase struct {
	footer  string
	visible bool

	styles config.Styles
}

func newDialogBase(styles config.Styles) dialogBase {
	return dialogBase{styles: styles}
}

func (b *dialogBase) SetFooter(footer string) {
	b.footer = footer
}

func (b *dialogBase) SetVisible(visible bool) {
	b.visible = visible
}

// render returns the dialog body, appending the footer as a body line when one
// is set. The footer is drawn with the dialog body text style so it sits on the
// frame background. The window manager draws the single frame and title around
// it.
func (b *dialogBase) render(body string) string {
	if b.footer == "" {
		return body
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, "", b.textStyle().Render(b.footer))
}

// textStyle returns the dialog body text style on the dialog frame background.
func (b *dialogBase) textStyle() lipgloss.Style {
	return b.styles.TextStyle.Inherit(b.styles.DialogBase())
}

// inputStyle returns the dialog input style on the dialog frame background.
func (b *dialogBase) inputStyle() lipgloss.Style {
	return b.styles.InputStyle.Inherit(b.styles.DialogBase())
}

// optionStyle returns the style for a Select option row. The focused row uses
// the injected dialog cursor style, falling back to the file panel cursor; the
// other rows use the injected dialog option style, falling back to the body
// text style. Both sit on the dialog frame background.
func (b *dialogBase) optionStyle(focused bool) lipgloss.Style {
	if focused {
		return b.styles.DialogCursorStyle.Inherit(b.styles.CursorStyle).Inherit(b.styles.DialogBase())
	}

	return b.styles.DialogOptionStyle.Inherit(b.styles.TextStyle).Inherit(b.styles.DialogBase())
}
