//nolint:testpackage // tests construct Model directly to reach addDialogWindow
package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/dialogs"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/wm"
)

func TestDialogWindowFitsContent(t *testing.T) {
	t.Parallel()

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
	}

	model.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	info := dialogs.NewInfo(model.styles)
	info.SetText("hi")
	info.SetVisible(true)

	id := model.addDialogWindow(info, "Info")

	win := model.main.WM().Get(id)
	if win == nil {
		t.Fatal("dialog window was not added")
	}

	wantW := max(lipgloss.Width(info.View()), wm.TitleWidth("Info")) + wm.FrameCols
	wantH := lipgloss.Height(info.View()) + wm.FrameRows

	if win.Width != wantW || win.Height != wantH {
		t.Errorf("dialog window = %dx%d, want %dx%d", win.Width, win.Height, wantW, wantH)
	}
}
