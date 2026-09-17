## 1. Select dialog cancel

- [x] 1.1 In `internal/ui/widgets/dialogs/select.go`, handle Escape by marking the dialog canceled, hiding it, and signaling `Done()`, and expose `Canceled()`
- [x] 1.2 In `internal/ui/dialogs.go`, make `Select` return no option when the dialog was canceled

## 2. Panel API

- [x] 2.1 In `internal/ui/widgets/panel/panel.go`, add `ColumnTitles()` and `SortBy(title)` (invert when the title is the current sort column, otherwise sort ascending by it), remove `CycleSort`, and share the re-sort and cursor-positioning logic

## 3. Menu window

- [x] 3.1 In `internal/ui/ui.go`, make the Menu action open the sort window for the focused panel and apply the chosen column on the event loop (a message), leaving the sort unchanged on cancel

## 4. Tests

- [x] 4.1 Dialog test: Escape cancels the Select dialog and reports no option
- [x] 4.2 Panel test: `SortBy` a different column sorts ascending, and `SortBy` the current column inverts its direction
- [x] 4.3 UI test: the Menu action opens the sort window (no mock) and applies the chosen column to the focused panel
- [x] 4.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 4.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 4.6 pty smoke (`make run`): F2 opens the column list, choosing a column sorts the panel, Escape cancels, and F10 quits

## 5. Modal dialogs

- [x] 5.1 In `internal/ui/widgets/wm/wm.go`, expose whether a dialog (non-panel) window currently holds focus, or add a modal flag the mainform can query
- [x] 5.2 In `internal/ui/widgets/mainform/mainwindow.go`, while a dialog is open, stop routing F1–F9 to the button bar and ignore mouse clicks outside the dialog; keep F10 quitting
- [x] 5.3 Test: with a dialog open, F1–F9 do not activate the bar, a click outside does not change focus, and F10 still quits
- [x] 5.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 5.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 5.6 pty smoke (`make run`): with the sort window open, F2 does not open a second window and a click outside is ignored, while F10 still quits
