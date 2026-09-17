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
// is set. The window manager draws the single frame and title around it.
func (b *dialogBase) render(body string) string {
	if b.footer == "" {
		return body
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, "", b.footer)
}
