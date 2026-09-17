## 1. Panel keymap: right arrow jumps to the bottom

- [x] 1.1 In `navigationKeyMap` (`internal/ui/widgets/panel/panel.go`), bind `right` to `keyMap.GotoBottom` instead of `keyMap.LineDown`, and update the function comment to say Left jumps to the first row and Right jumps to the last row
- [x] 1.2 Confirm `left` remains bound to `GotoTop` and the default Up/Down single-row stepping is untouched

## 2. Tests

- [x] 2.1 Replace `TestRightMovesDownAndClamps` with a test asserting a single Right moves the cursor from the first row to the last row of the listing
- [x] 2.2 Add a test asserting Right while already on the last row keeps the cursor on the last row
- [x] 2.3 Update `TestReloadResetsCursorToTop` and any other test that used repeated Rights to reach a specific row so it uses Down (or positions the cursor another way) instead of assuming Right steps one row
- [x] 2.4 Verify the Left-jumps-to-top test still passes and add/adjust coverage so Left/Right are a top/bottom pair
- [x] 2.5 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 2.6 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 2.7 Manual pty smoke (`make run`): in a directory with many entries, Right moves the cursor to the last row and Left returns it to the first; Down/Up still step one row; Enter/Backspace and the other panel are unaffected; F10 quits
