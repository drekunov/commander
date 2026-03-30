package wm

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Window wraps a tea.Model with position, size, z-index, and focus state.
type Window struct {
	ID      int
	Content tea.Model
	X, Y    int
	Width   int
	Height  int
	ZIndex  int
	Title   string
	Visible bool
	Focused bool

	// Drag state.
	dragging bool
	dragOffX int
	dragOffY int

	// Resize state.
	resizing    bool
	resizeOrigW int
	resizeOrigH int
	resizeOrigX int
	resizeOrigY int
}

// NewWindow creates a window with the given content, position, and size.
func NewWindow(id int, content tea.Model, x, y, w, h int) *Window {
	return &Window{
		ID:      id,
		Content: content,
		X:       x,
		Y:       y,
		Width:   w,
		Height:  h,
		Visible: true,
	}
}

// Rect returns the window's bounding rectangle.
func (w *Window) Rect() (x, y, width, height int) {
	return w.X, w.Y, w.Width, w.Height
}

// Contains returns true if the screen coordinate is inside the window.
func (w *Window) Contains(sx, sy int) bool {
	return sx >= w.X && sx < w.X+w.Width && sy >= w.Y && sy < w.Y+w.Height
}

// InTitleBar returns true if the coordinate is on the window's title bar (top row).
func (w *Window) InTitleBar(sx, sy int) bool {
	return sy == w.Y && sx >= w.X && sx < w.X+w.Width
}

// InResizeGrip returns true if the coordinate is on the bottom-right resize corner.
func (w *Window) InResizeGrip(sx, sy int) bool {
	return sx >= w.X+w.Width-2 && sy == w.Y+w.Height-1
}
