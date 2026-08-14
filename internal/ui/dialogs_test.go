//nolint:testpackage // tests inject the sendMsg seam and read unexported state
package ui

import (
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/drekunov/gc/internal/config"
	"github.com/drekunov/gc/internal/ui/widgets/mainform"
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
