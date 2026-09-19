## Why

After the user picks a column in the sort window opened by the Menu button, the panel beneath comes back with no visible selection cursor, and the Menu button stays stuck in its pressed style. Both defects are visible side effects of how dialog windows are added and removed outside the mainform update loop.

## What Changes

- Re-sync panel cursor visibility whenever a dialog window closes, so the panel that regains window focus renders its cursor again. Today the sync runs only while `mainform.Update` processes a message; the window manager removes a dialog from a goroutine, so no sync runs afterwards and the table stays blurred.
- Clear the function-button bar's pressed flash when the activation is dispatched, so the button returns to its normal style as soon as its action runs (including while the resulting dialog is open) instead of staying pressed until the next key reaches the bar.
- Add regression tests covering both behaviors.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `panel/navigation`: the "Only the focused panel shows the selection cursor" requirement is extended so cursor visibility is restored after a dialog window closes, not only while a mainform message is being processed.
- `ui/function-button-bar`: the "Activated buttons give visible feedback" requirement is clarified so the pressed style ends when the activation is dispatched rather than lingering for the lifetime of an opened dialog.

## Impact

- **Code**: `internal/ui/dialogs.go`, `internal/ui/ui.go`, `internal/ui/widgets/mainform/mainwindow.go`, `internal/ui/widgets/buttonbar/buttonbar.go`.
- **Tests**: new regression tests in `internal/ui/sort_test.go` / `internal/ui/dialogs_test.go` and `internal/ui/widgets/buttonbar/buttonbar_test.go`.
- **UI API**: internal-only; no public interface signature changes.
- **Dependencies**: none added.
