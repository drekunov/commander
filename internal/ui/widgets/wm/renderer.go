package wm

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/drekunov/gc/internal/config"
)

// renderer draws a window's frame, title bar, and resize grip using the
// configured theme. It is a pure state-to-string transform.
type renderer struct {
	styles config.Styles
}

// renderWindow draws a window frame with title bar and content.
func (r *renderer) renderWindow(win *Window) string {
	// Inner content area (subtract the frame columns and the title row).
	innerW := win.Width - FrameCols
	innerH := win.Height - FrameRows

	if innerW < 1 {
		innerW = 1
	}

	if innerH < 1 {
		innerH = 1
	}

	// Render the tea.Model content.
	content := win.Content.View()

	// Build title bar.
	titleBar := r.buildTitleBar(win, innerW)

	// Build resize grip.
	grip := r.buildResizeGrip(innerW)

	// Frame the content with a border — blue background fills the side borders.
	bgColor := r.styles.DialogBoxStyle.GetBackground()
	border := r.styles.DialogBoxStyle.GetBorderStyle()
	borderStyle := lipgloss.NewStyle().
		Foreground(r.styles.DialogBoxStyle.GetBorderBottomForeground()).
		Background(bgColor)
	styledLeft := borderStyle.Render(border.Left)
	styledRight := borderStyle.Render(border.Right)

	// Compose: title bar + content + grip, with side borders.
	var out strings.Builder

	out.WriteString(titleBar)
	out.WriteByte('\n')

	// Content lines with side borders, clipped/padded to exactly innerH lines
	// of innerW columns. Lipgloss Height/Width corrupt multi-line ANSI content
	// (it interleaves blank rows), so framing is done line by line.
	padStyle := lipgloss.NewStyle().Background(bgColor)

	contentLines := strings.Split(content, "\n")

	for row := range innerH {
		out.WriteString(styledLeft)

		line := ""
		if row < len(contentLines) {
			line = contentLines[row]
		}

		// Pad or clip line to innerW.
		lineW := lipgloss.Width(line)
		if lineW > innerW {
			line = ansi.Truncate(line, innerW, "")
		} else if lineW < innerW {
			line += padStyle.Render(strings.Repeat(" ", innerW-lineW))
		}

		out.WriteString(line)
		out.WriteString(styledRight)
		out.WriteByte('\n')
	}

	// Bottom border with resize grip.
	out.WriteString(grip)

	return out.String()
}

func (r *renderer) buildTitleBar(win *Window, innerW int) string {
	border := r.styles.DialogBoxStyle.GetBorderStyle()
	bgColor := r.styles.DialogBoxStyle.GetBackground()
	borderFg := r.styles.DialogBoxStyle.GetBorderBottomForeground()

	// NC-style title: [ Title ]
	title := "[ " + win.Title + " ]"

	titleLen := lipgloss.Width(title)
	if titleLen > innerW {
		title = title[:innerW]
		titleLen = innerW
	}

	padding := innerW - titleLen
	left := padding / 2
	right := padding - left

	borderStyle := lipgloss.NewStyle().Foreground(borderFg).Background(bgColor)

	var titleStyle lipgloss.Style
	if win.Focused {
		// Focused: bright white title, bold — matches NC active dialog style.
		titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(bgColor).
			Bold(true)
	} else {
		// Unfocused: same cyan as the border, no bold.
		titleStyle = lipgloss.NewStyle().
			Foreground(borderFg).
			Background(bgColor)
	}

	leftPart := border.TopLeft + strings.Repeat(border.Top, left)
	rightPart := strings.Repeat(border.Top, right) + border.TopRight

	return borderStyle.Render(leftPart) + titleStyle.Render(title) + borderStyle.Render(rightPart)
}

func (r *renderer) buildResizeGrip(innerW int) string {
	border := r.styles.DialogBoxStyle.GetBorderStyle()
	bgColor := r.styles.DialogBoxStyle.GetBackground()

	gripRune := "◢"
	bottomLen := max(innerW-1, 0)

	grip := border.BottomLeft +
		strings.Repeat(border.Bottom, bottomLen) +
		gripRune +
		border.BottomRight

	return lipgloss.NewStyle().
		Foreground(r.styles.DialogBoxStyle.GetBorderBottomForeground()).
		Background(bgColor).
		Render(grip)
}
