package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drekunov/gc/internal/app"
)

var errReadDir = errors.New("boom")

type fakeUI struct {
	infoCalls chan string
	data      chan []app.AttributeList
}

func (f *fakeUI) Info(_ context.Context, title, _, _ string) {
	if f.infoCalls != nil {
		f.infoCalls <- title
	}
}

func (f *fakeUI) Error(_, _ string) {}

func (f *fakeUI) Warning(_, _ string) {}

func (f *fakeUI) Confirm(_ context.Context, _, _ string) bool {
	return false
}

func (f *fakeUI) Input(_ context.Context, _, _, _ string) string {
	return ""
}

func (f *fakeUI) Password(_ context.Context, _, _ string) string {
	return ""
}

func (f *fakeUI) Select(_ context.Context, _, _ string, _ []string) string {
	return ""
}

func (f *fakeUI) SelectMultiple(_ context.Context, _, _ string, _ []string) []string {
	return nil
}

func (f *fakeUI) SetData(data []app.AttributeList) {
	if f.data != nil {
		f.data <- data
	}
}

type fakeConnector struct {
	rows []app.AttributeList
	err  error
}

func (f *fakeConnector) Name() string {
	return "Fake"
}

func (f *fakeConnector) ReadDir(string) ([]app.AttributeList, error) {
	if f.err != nil {
		return nil, f.err
	}

	return f.rows, nil
}

func (f *fakeConnector) ReadFile(string) ([]byte, error) {
	return nil, nil
}

func TestRunShowsDialogOnReadDirError(t *testing.T) {
	t.Parallel()

	uiFake := &fakeUI{infoCalls: make(chan string, 1)}
	connFake := &fakeConnector{err: errReadDir}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInstance := app.New(uiFake, connFake)

	errCh := make(chan error, 1)

	go func() {
		errCh <- appInstance.Run(ctx)
	}()

	select {
	case title := <-uiFake.infoCalls:
		if title != connFake.Name() {
			t.Errorf("info dialog title = %q, want connector name", title)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("info dialog was not shown on read error")
	}

	err := <-errCh
	if err == nil {
		t.Error("Run should return an error when ReadDir fails")
	}
}

func TestRunSetDataOnSuccess(t *testing.T) {
	t.Parallel()

	rows := []app.AttributeList{
		{{AttrName: "Name", AttrValue: "Name"}},
		{{AttrName: "Name", AttrValue: "a.txt"}},
	}

	uiFake := &fakeUI{data: make(chan []app.AttributeList, 1)}
	connFake := &fakeConnector{rows: rows}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appInstance := app.New(uiFake, connFake)

	errCh := make(chan error, 1)

	go func() {
		errCh <- appInstance.Run(ctx)
	}()

	select {
	case got := <-uiFake.data:
		if len(got) != 2 {
			t.Errorf("SetData received %d rows, want 2", len(got))
		}
	case <-time.After(2 * time.Second):
		t.Fatal("SetData was not called")
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after context cancellation")
	}
}
