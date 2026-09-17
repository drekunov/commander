## 1. Panel: remember the directory being left

- [x] 1.1 Add a `pendingSelect string` field to `panel.Model` in `internal/ui/widgets/panel/panel.go`, documented as the entry name to reselect on the next delivered listing
- [x] 1.2 In `ascendCmd`, after the root no-op check and before returning the `NavigateMsg`, set `m.pendingSelect = filepath.Base(m.dir)`
- [x] 1.3 In `navigationCmd`, clear `m.pendingSelect` when emitting a descend `NavigateMsg` so a stale ascent name cannot leak into the new listing

## 2. Panel: place the cursor when the listing arrives

- [x] 2.1 Add a helper that returns the display-row index for `m.pendingSelect`: match the `Name` attribute of each entry in `m.rows` (reusing `attrValue`/`entryName`) and add one when `showParent()` is true; return 0 when the name is empty or unmatched
- [x] 2.2 In `SetData`, replace `m.tableView.SetCursor(0)` with the computed index and clear `m.pendingSelect` after use
- [x] 2.3 Update the `SetData` doc comment to state that the cursor is restored to the just-left directory on ascent and reset to the first row otherwise

## 3. Tests and verification

- [x] 3.1 Panel test: after descending into a child and delivering the parent listing, the cursor is on the child's row (parent-row offset accounted for)
- [x] 3.2 Panel test: ascending via Enter on the ".." row restores the cursor the same way as Backspace
- [x] 3.3 Panel test: descend /a → b → c, then ascend one level at a time, asserting the cursor is on c in /a/b and then on b in /a
- [x] 3.4 Panel test: when the parent listing lacks the just-left entry, the cursor is on the first row
- [x] 3.5 Confirm descending (and a plain reload) still places the cursor on the first row
- [x] 3.6 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.7 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.8 Manual pty smoke (`make run`): descend into a directory, then Backspace — the cursor returns to the directory just left; repeat at a second level to confirm per-level memory; ascending to a parent that lacks the entry (e.g. after the child is gone) falls back to the first row; Enter still starts at the first row; the other panel and F10 quit are unaffected
