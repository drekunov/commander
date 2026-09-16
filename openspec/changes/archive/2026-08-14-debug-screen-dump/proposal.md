## Why

Testing the TUI (layout, window manager compositing, dialogs) requires seeing the exact rendered frame the user sees, but there is no way to capture it. A golden screen dump captured on a hotkey makes debugging and automated/visual testing possible without screen-scraping the terminal.

## What Changes

- Add a debug hotkey (F12) that captures the full composited screen (the same string `View()` renders) and writes it to a file in the process working directory.
- The dump file is timestamped (e.g. `screen-20260814-153045.ansi`) and written with raw bytes including ANSI styles, so it can be replayed with `cat` for a faithful reproduction of colors and layout.
- The dump is written asynchronously through a Bubble Tea command so the event loop is not blocked; the render is captured synchronously at keypress time, so the frame matches what is on screen when the key is pressed.
- A `Dumper` seam (file-write interface) is introduced so the dump path is unit-testable without touching the real filesystem.
- Dialogs and overlays are included in the dump automatically, since the captured string is the complete composited screen.

## Capabilities

### New Capabilities
- `ui/screen-dump`: Pressing the debug key writes the current full rendered screen to a timestamped file in the working directory without disrupting the running UI.

### Modified Capabilities
- None. This is a new capability; no existing `openspec/specs/` capability changes.

## Impact

- **Code**: `internal/ui/ui.go` (key handling + dump command), a new `internal/ui/dump.go` (dumper seam + filename generation), `cmd/commander/main.go` unchanged.
- **UI API**: `ui.Model` gains an unexported debug key handler; the `app.UI` interface is unchanged. F10 quit behavior is untouched.
- **Dependencies**: none added (stdlib `os`/`time` only).
- **Tests**: new unit tests for the dumper (content + filename) and for the key handler wiring a fake dumper; run with `go test -race ./...`.
- **Docs**: AGENTS.md key facts updated to mention F12 screen dump.
