package wm

import (
	"strings"

	"github.com/charmbracelet/x/ansi"
)

// canvas composites ANSI-styled text blocks onto a fixed-size grid.
// It works line-by-line, preserving ANSI escape sequences intact.
type canvas struct {
	lines         []string
	width, height int
}

func newCanvas(width, height int) *canvas {
	blank := strings.Repeat(" ", width)
	lines := make([]string, height)
	for i := range lines {
		lines[i] = blank
	}
	return &canvas{lines: lines, width: width, height: height}
}

// stamp overlays rendered text onto the canvas at position (ox, oy).
// Each line of the rendered text replaces the corresponding segment of
// the canvas line, using ANSI-aware width calculations.
func (c *canvas) stamp(ox, oy int, rendered string) {
	srcLines := strings.Split(rendered, "\n")
	for dy, srcLine := range srcLines {
		row := oy + dy
		if row < 0 || row >= c.height {
			continue
		}

		srcW := ansi.StringWidth(srcLine)
		if srcW == 0 || ox >= c.width {
			continue
		}

		// Clip source line to canvas bounds.
		clippedSrc := srcLine
		startCol := ox
		if startCol < 0 {
			// Trim left portion that falls off-screen.
			clippedSrc = ansi.TruncateLeft(clippedSrc, -startCol, "")
			srcW = ansi.StringWidth(clippedSrc)
			startCol = 0
		}

		endCol := startCol + srcW
		if endCol > c.width {
			clippedSrc = ansi.Truncate(clippedSrc, c.width-startCol, "")
			endCol = c.width
		}

		// Build: [left of canvas] + [clipped source] + [right of canvas]
		bg := c.lines[row]
		left := ansi.Truncate(bg, startCol, "")
		leftW := ansi.StringWidth(left)
		// Pad left if it's too short.
		if leftW < startCol {
			left += strings.Repeat(" ", startCol-leftW)
		}

		right := ansi.TruncateLeft(bg, endCol, "")

		c.lines[row] = left + clippedSrc + right
	}
}

// String renders the canvas into a single string.
func (c *canvas) String() string {
	return strings.Join(c.lines, "\n")
}

// Overlay composites foreground text on top of background text
// within the given dimensions.
func Overlay(background, foreground string, width, height int) string {
	c := newCanvas(width, height)
	c.stamp(0, 0, background)
	c.stamp(0, 0, foreground)
	return c.String()
}
