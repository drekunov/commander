//nolint:testpackage // tests inject the sendMsg seam and read unexported state
package ui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/buttonbar"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
	"github.com/drekunov/gc/internal/ui/widgets/topmenu"
)

func newTestModel() *Model {
	return &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
	}
}

func TestInfoReturnsOnCancelledContext(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})

	go func() {
		model.Info(ctx, "title", "footer", "message")
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Info did not return after context cancellation")
	}
}

func TestConfirmReturnsOnCancelledContext(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := make(chan bool, 1)

	go func() {
		result <- model.Confirm(ctx, "title", "message")
	}()

	select {
	case got := <-result:
		if got {
			t.Error("Confirm should return false on cancellation")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Confirm did not return after context cancellation")
	}
}

func TestInputReturnsOnCancelledContext(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := make(chan string, 1)

	go func() {
		result <- model.Input(ctx, "title", "footer", "message")
	}()

	select {
	case got := <-result:
		if got != "" {
			t.Errorf("Input = %q, want empty on cancellation", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Input did not return after context cancellation")
	}
}

func TestErrorReturnsOnQuit(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.quit = make(chan struct{})

	done := make(chan struct{})

	go func() {
		model.Error("title", "message")
		close(done)
	}()

	close(model.quit)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Error did not return after quit (dialog deadlock)")
	}
}

func TestRepeatedMockActivationReusesDialog(t *testing.T) {
	t.Parallel()

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
		quit:    make(chan struct{}),
	}

	defer model.stop()

	model.showMockText("Copy is not implemented yet")

	model.mockMu.Lock()
	first := model.mockDialog
	model.mockMu.Unlock()

	model.showMockText("Edit is not implemented yet")

	model.mockMu.Lock()
	second := model.mockDialog
	model.mockMu.Unlock()

	if first == nil || second == nil {
		t.Fatal("mock dialog was not created")
	}

	if first != second {
		t.Error("repeated activation opened a second mock dialog instead of reusing one")
	}

	if got := len(model.main.WM().Windows()); got != 3 {
		t.Errorf("window count = %d, want 3 (two panels and one mock dialog)", got)
	}
}

func TestMockDialogTitleIsUnimplemented(t *testing.T) {
	t.Parallel()

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
		quit:    make(chan struct{}),
	}

	defer model.stop()

	model.showMockText("View is not implemented yet")

	model.mockMu.Lock()
	id := model.mockWinID
	model.mockMu.Unlock()

	win := model.main.WM().Get(id)
	if win == nil {
		t.Fatal("mock dialog window was not created")
	}

	if win.Title != "Unimplemented" {
		t.Errorf("mock dialog title = %q, want Unimplemented", win.Title)
	}
}

// showMockText runs on the event loop; posting to the program from there would
// deadlock on the program's unbuffered, single-reader message channel.
func TestShowMockDoesNotPostToProgram(t *testing.T) {
	t.Parallel()

	posted := 0

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) { posted++ },
		quit:    make(chan struct{}),
	}

	defer model.stop()

	model.showMockText("View is not implemented yet")

	if posted != 0 {
		t.Errorf("showMock posted %d message(s) from the event loop; want 0", posted)
	}
}

func TestPullDnActivatesTopMenu(t *testing.T) {
	t.Parallel()

	model := newTestModel()

	model.handleBarActivation(buttonbar.ActionPullDn)

	if !model.main.TopMenu().Active() {
		t.Error("9PullDn did not activate the top menu")
	}

	if model.mockDialog != nil {
		t.Error("9PullDn opened a mock dialog instead of the top menu")
	}
}

func TestTopMenuFilesItemDispatchesBarAction(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.main.TopMenu().Activate()

	cmd := model.handleTopMenuActivation(topmenu.ActivateMsg{Action: buttonbar.ActionQuit})
	if cmd == nil {
		t.Fatal("a Files action produced no command")
	}

	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("Files action resolved to %T, want tea.QuitMsg", cmd())
	}

	if model.main.TopMenu().Active() {
		t.Error("the top menu stayed active after an item was activated")
	}
}

func TestTopMenuPlaceholderReportsLabel(t *testing.T) {
	t.Parallel()

	model := &Model{
		main:    mainform.New(config.Styles{}),
		styles:  config.Styles{},
		sendMsg: func(tea.Msg) {},
		quit:    make(chan struct{}),
	}

	defer model.stop()

	model.Update(topmenu.ActivateMsg{Label: "Find file"})

	model.mockMu.Lock()
	dialog := model.mockDialog
	model.mockMu.Unlock()

	if dialog == nil {
		t.Fatal("a placeholder item did not open the mock dialog")
	}

	if !strings.Contains(dialog.View(), "Find file") {
		t.Errorf("mock dialog does not name the placeholder item:\n%s", dialog.View())
	}
}
