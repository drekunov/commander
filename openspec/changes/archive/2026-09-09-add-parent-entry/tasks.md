## 1. Parent (..) row synthesis in the panel

- [x] 1.1 In `internal/ui/widgets/panel/panel.go` add a helper `hasParent(dir)` (true when `filepath.Dir(dir) != dir`) and track the current column titles from the delivered header row so a synthetic row can be built for any connector schema (fallback to a single `Name` column when no titles are known)
- [x] 1.2 Synthesize the `..` display row in `SetData`: cells derived from column titles (`Name`→`..`, `path`→the current dir, `IsDir`→`true`, others blank); when the current dir has a parent, render `[..] + entries`, otherwise render entries alone
- [x] 1.3 Empty non-root directories render only the `..` row (rows are cleared of stale entries, the `..` row remains); the filesystem root renders no `..` row
- [x] 1.4 Keep the real connector entries in a separate slice (`m.rows`) and map the table cursor to either the parent row (cursor 0 when a parent is present) or the real entry at `cursor-1`; route `selectedEntry`/selection helpers through this mapping

## 2. Enter on `..` and navigation semantics

- [x] 2.1 Extend `panel.navigationCmd`: Enter on the parent row emits the same `NavigateMsg` for `filepath.Dir(dir)` that Backspace emits; directory-descend and file-no-op branches are unchanged
- [x] 2.2 Confirm Backspace behavior is unchanged (ascends to the parent; no-op at the root) and that Enter at the root keeps its existing semantics (no `..` row to ascend)

## 3. Tests and verification

- [x] 3.1 Panel tests: a non-root listing leads with a `..` row and the root listing does not; Enter on `..` yields a `NavigateMsg` to the parent; the real directory below `..` still descends (index mapping); an empty non-root directory renders only `..`; Right/Left stepping works across the offset
- [x] 3.2 Update existing panel navigation tests whose cursor-0 assumption shifts when a parent row is present
- [x] 3.3 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.4 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.5 Manual pty smoke (`make run`): a non-root panel shows `..` as its first row; Enter on it returns to the parent; Backspace still ascends; the root panel shows no `..`; F10 quits; F12 still writes a screen dump
