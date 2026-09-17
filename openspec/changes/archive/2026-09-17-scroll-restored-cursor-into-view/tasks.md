## 1. Panel: scroll the selected row into view

- [x] 1.1 In `SetData` (`internal/ui/widgets/panel/panel.go`), replace the `m.tableView.SetCursor(m.cursorForPendingSelect())` call with a helper that positions the cursor via `m.tableView.GotoTop()` then `m.tableView.MoveDown(idx)`, where `idx` is the value from `cursorForPendingSelect`
- [x] 1.2 Update the `SetData` doc comment (and the helper's) to state that the selected row is scrolled into view, including the restored directory after ascent

## 2. Tests

- [x] 2.1 Panel test: build a listing taller than the panel viewport with a directory as the last entry; set width/height small; descend into that directory, ascend, then render `View()` and assert the restored entry's name appears (it does not with the old `SetCursor` behavior)
- [x] 2.2 Panel test: after a first-row reset (descend or plain reload) the first row is visible in the rendered `View()`
- [x] 2.3 Confirm the existing cursor-restoration tests (`TestBackspaceRestoresCursorToChild`, `TestEnterOnParentRowRestoresCursorToChild`, `TestAscendRestoresCursorAtEachLevel`, `TestAscendFallsBackWhenChildMissing`) still pass unchanged
- [x] 2.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 2.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 2.6 Manual pty smoke (`make run`): create a directory with more entries than fit the panel, scroll to a directory near the bottom, enter it, then Backspace — the restored row is highlighted and on screen; a plain reload/descend still shows the first row; arrows still scroll normally; F10 quits
