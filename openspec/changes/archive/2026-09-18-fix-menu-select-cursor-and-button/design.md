## Context

See proposal.md - Why. Existing constraints that shape the fix:

- Dialog calls in `internal/ui/dialogs.go` add a window with `m.addDialogWindow(...)` (which runs through the window manager) and then block on `select { <-ctx.Done(); <-m.quit; <-dialog.Done() }`.
- Opening a dialog is visible to the event loop: the call sends `tea.ResumeMsg{}` after adding the window, so `mainform.Update` processes it and calls `syncPanelFocus()`, blurring the panels.
- Closing a dialog is not visible to the event loop: the window is removed either by a deferred `m.main.WM().Remove(id)` in `dialogs.go` or by `waitMock`, both running off the Bubble Tea goroutine. `mainform.syncPanelFocus()` therefore never runs again, and the panel's bubbles table stays blurred even though `WM.Remove` re-focuses its window.
- `buttonbar.Model.pressed` is cleared only at the start of `buttonbar.Update`. `mainform.handleKey` never forwards keys to the bar while a dialog is open, so the pressed flash survives for the dialog's whole lifetime.
- `panel.tableStyles()` returns a blank `Selected` style whenever `tableView.Focused()` is false, so a blurred table renders no cursor.

## Goals / Non-Goals

**Goals:**

- The panel that regains focus after a dialog closes renders its selection cursor without requiring further input.
- A bar button returns to its normal style as soon as its activation is dispatched, and does not stay pressed while the resulting dialog is open.
- Keep the fix inside the existing message-driven focus-sync mechanism; do not add a focus interface to `tea.Model`.

**Non-Goals:**

- No change to `app.Dialog`/`app.Panel` interface signatures.
- No rework of the window manager's z-index or focus ordering.
- No change to the sort algorithm, the sort window's appearance, or other button behaviors.

## Decisions

### D1. Notify the event loop when a dialog window is removed

Centralize dialog removal in a helper on `ui.Model`:

```go
func (m *Model) closeDialogWindow(id int) {
    m.main.WM().Remove(id)
    m.sendMsg(tea.ResumeMsg{})
}
```

Replace the five `defer m.main.WM().Remove(id)` calls in `dialogs.go` with `defer m.closeDialogWindow(id)`, and replace the direct `m.main.WM().Remove(id)` in `waitMock` with `m.closeDialogWindow(id)`. The sent message lands in `mainform.Update`'s unhandled-message path, which already ends by calling `syncPanelFocus()`; the panel whose window `WM.Remove` re-focused is then un-blurred and its cursor reappears.

- *Alternative considered*: push focus into window content via a `Focusable` interface called from `WM.Add`/`Remove`. Rejected: `tea.Model` has no focus contract and this repeats the alternative already rejected by `fix-review-results` (D2); keeping the sync in mainform is localized.
- *Alternative considered*: call `m.main.SyncFocus()` directly after `Remove`. Rejected: `Remove` runs off the event loop and `panel.Focus()` mutates the bubbles table, which is not safe to touch concurrently with `Update`/`View`.
- *Alternative considered*: introduce a dedicated `mainform.FocusSyncMsg`. Rejected for now: the open path already uses `tea.ResumeMsg{}` for the same "wake the event loop to re-sync focus" purpose, and the unhandled-message path handles it without a new case; a dedicated message would be a larger refactor of both paths.

`Program.Send` is a no-op after the program terminates, so sending after `Remove` during quit cannot block or panic.

### D2. Clear the pressed flash when the activation is dispatched

Add `buttonbar.Model.ClearPressed()` (sets `pressed = 0`) and call it from `ui.Model.Update` in the `buttonbar.ActivateMsg` case before dispatching:

```go
case buttonbar.ActivateMsg:
    m.main.Bar().ClearPressed()
    return m, m.handleBarActivation(msg.Action)
```

The key/mouse message that activates the button still renders the pressed style for that frame, so the visible feedback requirement is met; the very next event-loop message (the `ActivateMsg`) clears it, so the button is unpressed while the dialog it opened is visible.

- *Alternative considered*: keep clearing in `buttonbar.Update` and forward every mainform message to the bar. Rejected: it would also clear the flash in the same `Update` that sets it (no visible feedback) unless carefully ordered, and it couples bar resets to unrelated messages.
- *Alternative considered*: clear the flash inside `mainform.handleKey` after `bar.Update`. Rejected: it runs in the same `Update` that set `pressed`, so `View` would never render the pressed frame.

### D3. Regression tests

- A ui-level test that opens the sort window with a stub dialog, closes it, and asserts the panel's table is focused (cursor visible) and the bar is not pressed.
- A buttonbar test that `ClearPressed()` returns the bar to its normal style after an activation.

## Risks / Trade-offs

- **`tea.ResumeMsg` is broadcast to all windows by `wm.Update`** → panels receive a no-op update. This already happens on dialog open; no new behavior.
- **Sort ordering**: `closeDialogWindow`'s `Send` is enqueued before the command's `sortColumnMsg` result, so focus is restored before the sort is applied. Even if ordering were reversed, the later resync would restore the cursor, so the final state is correct.
- **Clearing the pressed flash on `ActivateMsg` shortens the flash** → it is still rendered for the activating frame, which satisfies the "briefly" wording; verify with `make run`.

## Migration Plan

No migration. Apply in dependency order (buttonbar → ui/dialogs), keep the tree buildable, add tests, then run `go build ./...`, `go vet ./...`, `go test -race ./...`, and `make lint`.

## Open Questions

None.
