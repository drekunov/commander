package app_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/drekunov/gc/internal/app"
)

var errBoom = errors.New("boom")

const (
	attrPath  = "path"
	attrName  = "Name"
	attrIsDir = "IsDir"
)

type setDataCall struct {
	panel app.PanelID
	dir   string
	data  []app.AttributeList
}

type fakeUI struct {
	infoCh chan string
	errCh  chan string
	dataCh chan setDataCall
}

func newFakeUI() *fakeUI {
	return &fakeUI{
		infoCh: make(chan string, 16),
		errCh:  make(chan string, 16),
		dataCh: make(chan setDataCall, 16),
	}
}

func (f *fakeUI) Info(_ context.Context, title, _, _ string) {
	f.infoCh <- title
}

func (f *fakeUI) Error(title, _ string) {
	f.errCh <- title
}

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

func (f *fakeUI) SetData(panel app.PanelID, dir string, data []app.AttributeList) {
	f.dataCh <- setDataCall{panel: panel, dir: dir, data: data}
}

type fakeConnector struct {
	dirs map[string][]app.AttributeList
	errs map[string]error
}

func newFakeConnector() *fakeConnector {
	return &fakeConnector{
		dirs: map[string][]app.AttributeList{},
		errs: map[string]error{},
	}
}

func (f *fakeConnector) Name() string {
	return "Fake"
}

func (f *fakeConnector) ReadDir(path string) ([]app.AttributeList, error) {
	err := f.errs[path]
	if err != nil {
		return nil, err
	}

	return f.dirs[path], nil
}

func (f *fakeConnector) ReadFile(string) ([]byte, error) {
	return nil, nil
}

// runApp starts the app and cancels it when the test finishes.
func runApp(t *testing.T, appInstance *app.App) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)

	go func() {
		errCh <- appInstance.Run(ctx)
	}()

	t.Cleanup(func() {
		cancel()

		select {
		case err := <-errCh:
			if err != nil {
				t.Errorf("Run returned error: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Error("Run did not return after context cancellation")
		}
	})
}

func waitForData(t *testing.T, uiFake *fakeUI) setDataCall {
	t.Helper()

	select {
	case call := <-uiFake.dataCh:
		return call
	case <-time.After(2 * time.Second):
		t.Fatal("SetData was not called")

		return setDataCall{}
	}
}

func rows(names ...string) []app.AttributeList {
	header := app.AttributeList{
		{AttrName: attrPath, AttrValue: attrPath},
		{AttrName: attrName, AttrValue: attrName},
		{AttrName: attrIsDir, AttrValue: attrIsDir},
	}

	out := make([]app.AttributeList, 0, 1+len(names))
	out = append(out, header)

	for _, name := range names {
		out = append(out, app.AttributeList{
			{AttrName: attrPath, AttrValue: "/"},
			{AttrName: attrName, AttrValue: name},
			{AttrName: attrIsDir, AttrValue: false},
		})
	}

	return out
}

func TestRunDeliversRootToBothPanels(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.dirs["/"] = rows("etc", "tmp")

	runApp(t, app.New(uiFake, connFake))

	left := waitForData(t, uiFake)
	if left.panel != app.PanelLeft || left.dir != "/" {
		t.Errorf("first delivery = %v %q, want left panel at /", left.panel, left.dir)
	}

	right := waitForData(t, uiFake)
	if right.panel != app.PanelRight || right.dir != "/" {
		t.Errorf("second delivery = %v %q, want right panel at /", right.panel, right.dir)
	}
}

func TestRunStartupReadErrorShowsInfoAndKeepsRunning(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.errs["/"] = errBoom

	runApp(t, app.New(uiFake, connFake))

	for range 2 {
		select {
		case <-uiFake.infoCh:
		case <-time.After(2 * time.Second):
			t.Fatal("info dialog was not shown for the startup read error")
		}
	}

	select {
	case call := <-uiFake.dataCh:
		t.Errorf("SetData was called after startup errors: %+v", call)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestNavigateDeliversToRequestedPanel(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.dirs["/"] = rows("etc")
	connFake.dirs["/etc"] = rows("hosts")

	appInstance := app.New(uiFake, connFake)
	runApp(t, appInstance)

	waitForData(t, uiFake) // left
	waitForData(t, uiFake) // right

	appInstance.Navigate(app.PanelRight, "/etc")

	call := waitForData(t, uiFake)
	if call.panel != app.PanelRight {
		t.Errorf("navigate delivered to panel %v, want right", call.panel)
	}

	if call.dir != "/etc" {
		t.Errorf("navigate delivered dir %q, want /etc", call.dir)
	}

	if len(call.data) != 2 {
		t.Errorf("navigate delivered %d rows, want 2", len(call.data))
	}
}

func TestNavigateKeepsLatestRequest(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.dirs["/"] = rows("etc")
	connFake.dirs["/a"] = rows("a")
	connFake.dirs["/b"] = rows("b")
	connFake.dirs["/c"] = rows("c")

	appInstance := app.New(uiFake, connFake)
	runApp(t, appInstance)

	waitForData(t, uiFake) // left
	waitForData(t, uiFake) // right

	for _, dir := range []string{"/a", "/b", "/c"} {
		appInstance.Navigate(app.PanelRight, dir)
	}

	deadline := time.After(2 * time.Second)

	for {
		select {
		case call := <-uiFake.dataCh:
			if call.panel == app.PanelRight && call.dir == "/c" {
				return
			}
		case <-deadline:
			t.Fatal("latest navigation request /c was not delivered")
		}
	}
}

func TestNavigateErrorShowsDialogAndKeepsListing(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.dirs["/"] = rows("etc")
	connFake.errs["/secret"] = errBoom

	appInstance := app.New(uiFake, connFake)
	runApp(t, appInstance)

	waitForData(t, uiFake)
	waitForData(t, uiFake)

	appInstance.Navigate(app.PanelRight, "/secret")

	select {
	case <-uiFake.errCh:
	case <-time.After(2 * time.Second):
		t.Fatal("error dialog was not shown for an unreadable directory")
	}

	select {
	case call := <-uiFake.dataCh:
		t.Errorf("SetData was called after a failed navigation: %+v", call)
	case <-time.After(200 * time.Millisecond):
	}
}

func TestNavigateErrorThenCancelReturns(t *testing.T) {
	t.Parallel()

	uiFake := newFakeUI()
	connFake := newFakeConnector()
	connFake.dirs["/"] = rows("etc")
	connFake.errs["/secret"] = errBoom

	appInstance := app.New(uiFake, connFake)

	ctx, cancel := context.WithCancel(context.Background())

	errCh := make(chan error, 1)

	go func() {
		errCh <- appInstance.Run(ctx)
	}()

	waitForData(t, uiFake) // left
	waitForData(t, uiFake) // right

	appInstance.Navigate(app.PanelRight, "/secret")

	select {
	case <-uiFake.errCh:
	case <-time.After(2 * time.Second):
		cancel()
		t.Fatal("error dialog was not shown for an unreadable directory")
	}

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Errorf("Run returned error: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after cancellation")
	}
}
