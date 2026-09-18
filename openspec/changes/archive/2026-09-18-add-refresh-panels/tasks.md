## 1. App batch read

- [x] 1.1 Change `App.pending` (`internal/app/app.go`) from a single request to a slice and make `takePending` return it.
- [x] 1.2 Keep `Navigate` storing a single-element latest-wins batch, and add `Refresh(requests []NavRequest)` that stores the whole batch and wakes the run loop.
- [x] 1.3 Handle every request the run loop drains (one `handleNav` per request), preserving the existing startup and error behavior.

## 2. Panel selection preservation

- [x] 2.1 Add `panel.Model.PreserveSelection()` that stores the currently selected entry in `pendingSelect`, and call it from `ui.refreshPanels` for each panel before requesting the re-read.
- [x] 2.2 Verify a missing entry falls back to the first row and that navigating to a different directory still resets the cursor.

## 3. Mainform and ui key handling

- [x] 3.1 Export the modal state from `mainform` (e.g. `DialogOpen()`), and add `mainform.Panels()` returning the left and right panels in order.
- [x] 3.2 Add `panel.Model.ID() app.PanelID`.
- [x] 3.3 Add a `refresher func([]app.NavRequest)` field and `SetRefresher` to `ui.Model`, and handle `tea.KeyCtrlR` in `ui.Model.Update`: when no dialog is open, build the batch from both panels' non-empty directories and call the refresher.

## 4. Wiring

- [x] 4.1 Call `uiInstance.SetRefresher(appInstance.Refresh)` in `cmd/commander/main.go`.

## 5. Tests

- [x] 5.1 App test: `Refresh` with two requests delivers listings to both panels; `Navigate` latest-wins still passes.
- [x] 5.2 Panel test: re-delivering the same directory keeps the selected entry; removing it falls back to the first row.
- [x] 5.3 ui/mainform test: `Ctrl+R` while a dialog is open does nothing; with no dialog it invokes the refresher with both panels.

## 6. Verification

- [x] 6.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 6.2 Run `make lint` and confirm golangci-lint passes.
- [ ] 6.3 Manual smoke via `make run`: create a file in a panel's directory from another shell, press `Ctrl+R`, and confirm both panels refresh with the selection kept; press `Ctrl+R` with a dialog open and confirm nothing changes.
