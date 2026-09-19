## 1. Focus resync when a dialog closes

- [x] 1.1 Add a `closeDialogWindow(id int)` helper on `ui.Model` (`internal/ui/dialogs.go`) that calls `m.main.WM().Remove(id)` and then sends `tea.ResumeMsg{}` so the event loop re-runs the focus sync.
- [x] 1.2 Replace the five `defer m.main.WM().Remove(id)` calls in `dialogs.go` (`Info`, `runInput`, `Select`, `SelectMultiple`, `Confirm`) with `defer m.closeDialogWindow(id)`.
- [x] 1.3 Replace the direct `m.main.WM().Remove(id)` in `waitMock` (`internal/ui/ui.go`) with `m.closeDialogWindow(id)`.

## 2. Clear the pressed flash on dispatch

- [x] 2.1 Add an exported `ClearPressed()` method to `buttonbar.Model` (`internal/ui/widgets/buttonbar/buttonbar.go`) that resets `pressed` to 0.
- [x] 2.2 Call `m.main.Bar().ClearPressed()` in `ui.Model.Update`'s `buttonbar.ActivateMsg` case before `handleBarActivation`, keeping the pressed style for the activating frame.

## 3. Regression tests

- [x] 3.1 Add a ui-level test that opens the sort window with a stub/injected dialog, closes it, and asserts the focused panel's table is focused (cursor visible) afterward.
- [x] 3.2 Add a buttonbar test asserting `ClearPressed()` returns a pressed button to its normal style.

## 4. Verification

- [x] 4.1 Run `go build ./...`, `go vet ./...`, and `go test -race ./...` and keep them green.
- [x] 4.2 Run `make lint` (rewrites in place) and confirm golangci-lint passes.
- [ ] 4.3 Manual smoke via `make run`: activate Menu, choose a column, and confirm the panel cursor is visible and the Menu button is not pressed; open and close a mock dialog (e.g. F3 View) and confirm the cursor returns.
