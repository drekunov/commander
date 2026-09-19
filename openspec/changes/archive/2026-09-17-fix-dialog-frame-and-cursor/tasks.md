## 1. Dialog frame

- [x] 1.1 In `internal/ui/widgets/dialogs/base.go`, make `render` return the body with the footer appended as a line, and remove the box/title/placement and the now-unused border helpers
- [x] 1.2 In `internal/ui/dialogs.go`, size `addDialogWindow` to the wider of the content and the title plus the frame (clamped to the canvas) and center it

## 2. Cursor style

- [x] 2.1 In `internal/ui/widgets/dialogs/select.go`, render the focused option row with `styles.CursorStyle`

## 3. Tests

- [x] 3.1 Dialog test: `render` output has no border and includes the footer as a body line
- [x] 3.2 Dialog test: the Select dialog's focused row uses the cursor style
- [x] 3.3 UI test: `addDialogWindow` sizes the window to the content plus the frame
- [x] 3.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.6 pty smoke (`make run`): the sort window shows a single frame sized to its content, the focused row uses the cursor style, and F10 quits

## 4. Sort window caption

- [x] 4.1 In `internal/ui/ui.go`, open the sort window with the caption `Sort: Select a column` and no body prompt
- [x] 4.2 UI test: the sort window's caption is `Sort: Select a column` and it has no body prompt line
- [x] 4.3 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 4.4 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 4.5 pty smoke (`make run`): the sort window caption reads `Sort: Select a column` (not truncated) with no body prompt, and F10 quits
