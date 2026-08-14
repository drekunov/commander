package dialogs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
)

const (
	dialogFrameCols = 2
)

// dialogBase provides the shared title-bar/footer/border shell used by every
// dialog widget in this package.
type dialogBase struct {
	title   string
	footer  string
	width   int
	height  int
	visible bool

	styles config.Styles
}

func newDialogBase(styles config.Styles) dialogBase {
	return dialogBase{styles: styles}
}

func (b *dialogBase) SetTitle(title string) {
	b.title = title
}

func (b *dialogBase) SetFooter(footer string) {
	b.footer = footer
}

func (b *dialogBase) SetVisible(visible bool) {
	b.visible = visible
}

func (b *dialogBase) SetWidth(width int) {
	b.width = width
}

func (b *dialogBase) SetHeight(height int) {
	b.height = height
}

// render frames body content with the dialog box style, title bar, and footer.
func (b *dialogBase) render(body string) string {
	output := b.styles.DialogBoxStyle.
		Width(b.width).
		Height(b.height).
		Align(lipgloss.Center, lipgloss.Center).
		BorderTop(false).
		BorderBottom(false).
		Render(body)

	header := b.headerView(lipgloss.Width(output) - dialogFrameCols)
	footer := b.footerView(lipgloss.Width(output) - dialogFrameCols)

	output = lipgloss.JoinVertical(lipgloss.Center, header, output, footer)

	return lipgloss.Place(b.width, b.height, 0.2, 0.2, output)
}

func (b *dialogBase) headerView(width int) string {
	borderStyle := b.styles.DialogBoxStyle.GetBorderStyle()

	header := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		" "+b.title+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top),
	)

	header = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.TopLeft, header, borderStyle.TopRight)

	header = lipgloss.NewStyle().
		Foreground(b.styles.DialogBoxStyle.GetBorderBottomForeground()).
		Render(header)

	return header
}

func (b *dialogBase) footerView(width int) string {
	borderStyle := b.styles.DialogBoxStyle.GetBorderStyle()

	footer := lipgloss.PlaceHorizontal(
		width,
		lipgloss.Center,
		" "+b.footer+" ",
		lipgloss.WithWhitespaceChars(borderStyle.Top),
	)

	footer = lipgloss.JoinHorizontal(lipgloss.Center, borderStyle.BottomLeft, footer, borderStyle.BottomRight)

	footer = lipgloss.NewStyle().
		Foreground(b.styles.DialogBoxStyle.GetBorderBottomForeground()).
		Render(footer)

	return footer
}
