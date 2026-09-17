## 1. Dialog quit safety (D1)

- [x] 1.1 Add `case <-m.quit:` to the blocking `select` in `Info` (`internal/ui/dialogs.go`), returning immediately.
- [x] 1.2 Add `case <-m.quit:` to `runInput`, `Select`, `SelectMultiple`, and `Confirm`, returning each method's zero value.
- [x] 1.3 Confirm `Error`/`Warning` now return on quit via `m.quit` without signature changes.

## 2. Focus wiring (D2)

- [x] 2.1 Remove the unconditional `right.Focus()` from `mainform.New`.
- [x] 2.2 Add `mainform.syncPanelFocus()` that mirrors each panel window's `Focused` flag to `panel.Focus()`/`panel.Blur()`.
- [x] 2.3 Call `syncPanelFocus()` after construction and after `m.wm.Update(msg)` in `mainform.Update`.

## 3. Bar background, styling, and pressed flash (D3, D4, D7)

- [x] 3.1 Derive the per-segment background from `ButtonBarStyle` so the full-width strip is themed.
- [x] 3.2 Replace the hand-rolled `colorize` in `buttonbar.go` with lipgloss `Style.Copy()`/`Inherit` composition.
- [x] 3.3 Move the `pressed = 0` reset from `buttonbar.View()` to the start of `buttonbar.Update`.
- [x] 3.4 Update `buttonbar_test.go` for the new reset timing and themed background expectations.

## 4. Mock dialog coalescing (D5)

- [x] 4.1 Track the open mock dialog and its window id on `ui.Model`.
- [x] 4.2 On activation, update the existing mock dialog text in place when one is open; otherwise create it.
- [x] 4.3 Clear the tracked reference when the mock goroutine finishes (quit or `Done()`).

## 5. Reliable navigation (D6)

- [x] 5.1 Replace the `select`/`default` drop in `App.Navigate` with a mutex-protected latest-request slot plus a signal channel.
- [x] 5.2 Drain the latest-request slot in `App.Run` under the mutex.
- [x] 5.3 Verify the event loop never blocks and the newest request is always honored.

## 6. Panel matching, complexity, and accessors (D8, D9, D10)

- [x] 6.1 Add a shared case-insensitive `attrValue(entry, name)` helper in `panel.go` and use it in `entryName`, `entryIsDir`, and `parentRow`.
- [x] 6.2 Extract listing→columns/rows construction into `columnsFor`/`displayRowsFor` helpers, keeping the `colTitles` fallback with a clarifying comment.
- [x] 6.3 Make `Action.Label()`/`Name()` return `""` for out-of-range values instead of panicking.
- [x] 6.4 Return early after `applyData` in `mainform.Update` so `DataMsg` is not broadcast to every window.

## 7. Regression tests

- [x] 7.1 Add a test that starts `App.Run`, triggers a failed `Navigate`, then cancels the app context and asserts `Run` returns (no deadlock) within a timeout.
- [x] 7.2 Add a `panel` test asserting only the focused panel renders a cursor after focus changes.
- [x] 7.3 Add a `buttonbar` test asserting the full-width background is rendered and repeated activation reuses one mock dialog.
- [x] 7.4 Rename the `asModel` test helper and remove the hardcoded `Y:29` coupling in `mainwindow_test.go`.

## 8. Verification

- [x] 8.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 8.2 Run `make lint` (rewrites in place) and confirm golangci-lint passes.
- [x] 8.3 Manual smoke via `make run`: quit with an error dialog open exits cleanly; only the focused panel shows a cursor; the bar strip background is visible.
