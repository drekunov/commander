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

// String renders the canvas into a single string.
func (cvs *canvas) String() string {
	return strings.Join(cvs.lines, "\n")
}

// stamp overlays rendered text onto the canvas at position (ox, oy).
// Each line of the rendered text replaces the corresponding segment of
// the canvas line, using ANSI-aware width calculations.
func (cvs *canvas) stamp(offsetX, offsetY int, rendered string) {
	srcLines := strings.Split(rendered, "\n")

	for dy, srcLine := range srcLines {
		row := offsetY + dy
		if row < 0 || row >= cvs.height {
			continue
		}

		srcW := ansi.StringWidth(srcLine)
		if srcW == 0 || offsetX >= cvs.width {
			continue
		}

		// Clip source line to canvas bounds.
		clippedSrc := srcLine

		startCol := offsetX
		if startCol < 0 {
			// Trim left portion that falls off-screen.
			clippedSrc = ansi.TruncateLeft(clippedSrc, -startCol, "")
			srcW = ansi.StringWidth(clippedSrc)
			startCol = 0
		}

		endCol := startCol + srcW
		if endCol > cvs.width {
			clippedSrc = ansi.Truncate(clippedSrc, cvs.width-startCol, "")
			endCol = cvs.width
		}

		// Build: [left of canvas] + [clipped source] + [right of canvas]
		background := cvs.lines[row]
		left := ansi.Truncate(background, startCol, "")

		leftW := ansi.StringWidth(left)
		// Pad left if it's too short.
		if leftW < startCol {
			left += strings.Repeat(" ", startCol-leftW)
		}

		right := ansi.TruncateLeft(background, endCol, "")

		cvs.lines[row] = left + clippedSrc + right
	}
}

func newCanvas(width, height int) *canvas {
	blank := strings.Repeat(" ", width)

	lines := make([]string, height)
	for i := range lines {
		lines[i] = blank
	}

	return &canvas{lines: lines, width: width, height: height}
}

// Overlay composites foreground text on top of background text
// within the given dimensions.
func Overlay(background, foreground string, width, height int) string {
	cvs := newCanvas(width, height)
	cvs.stamp(0, 0, background)
	cvs.stamp(0, 0, foreground)

	return cvs.String()
}
