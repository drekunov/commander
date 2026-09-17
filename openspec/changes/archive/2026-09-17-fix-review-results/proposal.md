## Why

A comprehensive code review of the function-button-bar + panel-navigation work surfaced a blocking deadlock (quit while an error dialog is open hangs the process), a double-cursor focus regression, dead theming config, and a set of smaller correctness, robustness, and complexity issues. This change fixes the whole list so the work can land without regressions.

## What Changes

- **Fix the quit-while-dialog-open deadlock**: `Error`/`Warning` open dialogs with `context.Background()`, so an open error dialog never returns on F10-quit and the process hangs. Every dialog call now returns on quit regardless of the context passed.
- **Fix the focus regression**: `right.Focus()` at construction makes both panels render a selection cursor. Window focus is wired to panel cursor visibility so only the focused panel shows the cursor.
- **Render the bar's themed background**: `ButtonBarStyle` (`.button-bar`) is loaded but never used; the bottom strip's background is dead config. The bar now paints its full-width background and overlays the buttons.
- **Coalesce mock dialogs**: repeated function-key presses no longer stack an unbounded number of mock Info dialogs and goroutines.
- **Make navigation reliable**: a valid navigation request is no longer silently dropped when the request queue is full.
- **Remove state mutation from `View()`**: the button "pressed" flash is cleared in `Update`, not in `View`.
- **Harden the button API**: `Action.Label()`/`Name()` no longer panic on out-of-range values.
- **Make attribute matching consistent**: parent-row synthesis and entry classification use the same (case-insensitive) matching rules.
- **Reduce panel complexity**: extract the listing→rows rendering into a focused helper and trim the speculative fallback paths.
- **Fix `colorize` to preserve full style fidelity** instead of hand-copying a subset of lipgloss attributes.
- **Stop double-delivering `DataMsg`** through the window manager after it is already applied.
- **Add regression tests** for the deadlock and tighten test naming/coupling nits.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `ui/dialogs`: the "Dialogs never hang the application" requirement is extended so every dialog call, including `Error`/`Warning` opened with a non-cancellable context, returns on quit.
- `ui/function-button-bar`: the bar renders its full-width themed background and mock activation does not stack dialogs.
- `panel/navigation`: only the focused panel renders its selection cursor, and a valid navigation request is not silently dropped.

## Impact

- **Code**: `internal/ui/dialogs.go`, `internal/ui/ui.go`, `internal/app/app.go`, `internal/ui/widgets/mainform/mainwindow.go`, `internal/ui/widgets/panel/panel.go`, `internal/ui/widgets/buttonbar/buttonbar.go`, `internal/config/config.go`, `internal/config/styles.css`.
- **Tests**: new regression tests in `internal/ui/dialogs_test.go` and `internal/app/app_test.go`; adjustments to `internal/ui/widgets/panel/panel_test.go` and `internal/ui/widgets/mainform/mainwindow_test.go`.
- **UI API**: `Dialog.Error`/`Warning` behavior changes (they now honor app quit); the `Dialog` interface signature itself is unchanged unless a cancellable context must be threaded (internal-only callers).
- **Dependencies**: none added.
