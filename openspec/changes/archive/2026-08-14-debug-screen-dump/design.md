## Context

`ui.Model.Update` already owns global key handling (`internal/ui/ui.go` — F10 quits). `Model.View()` returns the full composited frame — all windows, dialogs, overlays, and ANSI styling — via `mainform.Model.View()` → `wm.Manager.View()`. `wm.View()` is a read-only, mutex-guarded snapshot. See proposal.md - Why for motivation.

Constraints shaping the approach:

- The event loop is Bubble Tea; anything slow in `Update` blocks rendering and input.
- `ui.Model` has the `program` handle; the `app.UI` interface must stay unchanged (this feature is UI-internal).
- The working directory is where cwd-style overrides already live; dumps belong there too.
- No logger exists in `ui`; the alt screen means we cannot print errors without corrupting the display.

## Goals / Non-Goals

**Goals:**

- Capture the exact frame being displayed on F12 and persist it to `screen-<timestamp>.ansi` in the working directory.
- Keep the write off the event loop so the UI never stutters.
- Make the dump path injectable so the behavior is unit-testable (spec `ui/screen-dump`).
- Land green: `go build`, `go test -race ./...`, `make lint`.

**Non-Goals:**

- No changes to `cmd/wm-example` (its F10 handler stays; the demo is not the app).
- No in-UI user feedback (toast/status line) for a successful or failed dump.
- No terminal-screen-scrape or frame-diff tooling — this change only produces the capture.
- No config/flag to relocate the dump directory; the working directory is fixed for now.

## Decisions

### D1. F12 is the debug key

F10 quits; F12 is conventionally a debug key and is currently unbound. The key is a named constant (`dumpScreenKey = tea.KeyF12`), so re-binding later is a one-line change. Ctrl+D was rejected because terminals commonly use it for EOF and it risks conflicting with line discipline.

### D2. Capture synchronously, write asynchronously

In `Model.Update`, on F12 we capture `m.View()` into a local value and return a `tea.Cmd` whose closure writes that captured string. The write therefore never blocks `Update` (spec — "Dump does not block the UI"), and because the frame is captured synchronously at keypress time, the dumped string is exactly what the current state renders (spec — "Frame matches what is displayed"). Dialogs/overlays come along for free since we capture the whole composited frame.

- *Alternative considered*: writing synchronously in `Update`. Rejected: a slow disk could freeze the TUI.
- *Alternative considered*: a separate goroutine writing to a channel. Rejected: Bubble Tea's `tea.Cmd` is the idiomatic async mechanism and is already unit-testable.

### D3. Dumper seam for testability

`internal/ui/dump.go` defines a tiny interface and a production implementation:

```go
type dumper interface {
    Dump(frame string) error
}

type fileDumper struct {
    now func() time.Time   // injectable clock, defaults to time.Now
}
```

`Model` gains an unexported `dumper` field, defaulted to `&fileDumper{}` in `New`. Tests in `package ui` replace it with a fake that records the frame. No new dependencies; stdlib `os`/`time` only.

### D4. Filename: `screen-<YYYYMMDD-HHMMSS>.ansi`, collision-suffixed

Filename generation is a pure function `dumpFilename(now time.Time) string` → `screen-20060102-150405.ansi`, and the timestamp encodes local capture time at second precision (spec — "Dump file naming"). To satisfy "successive dumps never overwrite each other" even for two presses within the same second, the writer checks existence and appends `-1`, `-2`, … until the name is free.

- *Alternative considered*: nanosecond timestamps. Rejected: they embed more precision than the spec's second-precision contract and make filenames ugly for humans to read during debugging.
- *Alternative considered*: single fixed name `screendump.txt`. Rejected: overwrites the previous capture, defeating history.

### D5. Failures are swallowed, not displayed

`Dump` errors are ignored by the key handler. Printing to stdout/stderr would corrupt the alt screen, and a dialog would be obnoxious. The absence of the file on disk is the signal. This satisfies spec — "Dump failure does not crash the application".

## Risks / Trade-offs

- **Frame captured may lag one render** if a keypress and a data-update land in the same frame → the dumped frame is the state at keypress time; this is the most useful definition for debugging and matches the spec. No mitigation needed.
- **Second-precision names can collide** → collision suffix (`-1`, `-2`) guarantees distinct files.
- **Silent failures are invisible** → acceptable for a debug tool; the file's absence is the diagnostic.

## Migration Plan

1. Add `internal/ui/dump.go` (`dumper` interface, `fileDumper`, `dumpFilename`) with unit tests.
2. Wire F12 handling into `Model.Update` with a `dumper` field defaulted in `New`.
3. Add key-handler tests with a fake dumper.
4. Verify `go build`, `go test -race ./...`, `make lint`; manual smoke with `make run` (F12 → dump file appears, colors intact under `cat`).

Rollback: revert the single commit; no interface or contract changes touch other packages.

## Open Questions

None.
