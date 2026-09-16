//nolint:testpackage // tests exercise the unexported dumper and key handler
package ui

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type recordingDumper struct {
	frames []string
}

func (r *recordingDumper) Dump(frame string) error {
	r.frames = append(r.frames, frame)

	return nil
}

func fixedClock(now time.Time) func() time.Time {
	return func() time.Time { return now }
}

func TestDumpFilename(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 14, 15, 30, 45, 0, time.UTC)

	want := "screen-20260814-153045.ansi"
	if got := dumpFilename(now); got != want {
		t.Errorf("dumpFilename(%v) = %q, want %q", now, got, want)
	}
}

func TestFileDumperWritesRawFrame(t *testing.T) {
	t.Parallel()

	frame := "\x1b[31mred\x1b[0m"

	dumper := &fileDumper{
		dir: t.TempDir(),
		now: fixedClock(time.Date(2026, 8, 14, 15, 30, 45, 0, time.UTC)),
	}

	err := dumper.Dump(frame)
	if err != nil {
		t.Fatalf("Dump() error = %v", err)
	}

	assertFileContent(t, filepath.Join(dumper.dir, "screen-20260814-153045.ansi"), frame)
}

func TestFileDumperCollisionSuffixes(t *testing.T) {
	t.Parallel()

	dumper := &fileDumper{
		dir: t.TempDir(),
		now: fixedClock(time.Date(2026, 8, 14, 15, 30, 45, 0, time.UTC)),
	}

	err := dumper.Dump("first")
	if err != nil {
		t.Fatalf("first Dump() error = %v", err)
	}

	err = dumper.Dump("second")
	if err != nil {
		t.Fatalf("second Dump() error = %v", err)
	}

	assertFileContent(t, filepath.Join(dumper.dir, "screen-20260814-153045.ansi"), "first")
	assertFileContent(t, filepath.Join(dumper.dir, "screen-20260814-153045-1.ansi"), "second")
}

func TestFileDumperMissingDirReturnsError(t *testing.T) {
	t.Parallel()

	dumper := &fileDumper{
		dir: filepath.Join(t.TempDir(), "missing"),
		now: fixedClock(time.Date(2026, 8, 14, 15, 30, 45, 0, time.UTC)),
	}

	err := dumper.Dump("frame")
	if err == nil {
		t.Error("Dump() error = nil, want an error for a missing directory")
	}
}

func TestModelDumpOnF12(t *testing.T) {
	t.Parallel()

	model := newTestModel()
	model.dumper = &recordingDumper{}

	model.main, _ = model.main.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	want := model.View()

	got, cmd := model.Update(tea.KeyMsg{Type: tea.KeyF12})
	if cmd == nil {
		t.Fatal("F12 returned a nil command")
	}

	if _, ok := got.(*Model); !ok {
		t.Fatalf("Update() returned %T, want *Model", got)
	}

	if msg := cmd(); msg != nil {
		t.Errorf("dump command returned %v, want nil", msg)
	}

	rec, ok := model.dumper.(*recordingDumper)
	if !ok {
		t.Fatalf("dumper is %T, want *recordingDumper", model.dumper)
	}

	if len(rec.frames) != 1 {
		t.Fatalf("dumper received %d frames, want 1", len(rec.frames))
	}

	if rec.frames[0] != want {
		t.Error("dumped frame does not equal View() at keypress time")
	}
}

func assertFileContent(t *testing.T, path, want string) {
	t.Helper()

	got, err := os.ReadFile(path) //nolint: gosec
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}

	if string(got) != want {
		t.Errorf("%s = %q, want %q", path, got, want)
	}
}
