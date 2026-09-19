## 1. Sortable size

- [x] 1.1 In `internal/app/connector.go`, add a `Size` type (an `int64` whose `String()` renders the human-readable size) and move the human-size formatting there
- [x] 1.2 In `internal/connectors/filesystem/connector.go`, return the sortable `app.Size` for the Size column while keeping `<DIR>` for directories

## 2. Panel sorting

- [x] 2.1 In `internal/ui/widgets/panel/panel.go`, add the sort state (`sortCol`, `sortAsc`, default Name ascending) and a `CycleSort` method that advances the cycle
- [x] 2.2 Add a `sortRows` helper that orders directories above files and each group by the active column and direction, comparing `app.Size` numerically and other values as case-insensitive strings
- [x] 2.3 Apply the sort in `SetData` and keep the selected entry selected across a re-sort (reuse the pending-select mechanism)
- [x] 2.4 Render the active column's header with a direction marker in `applyColumnWidths`

## 3. Routing

- [x] 3.1 In `internal/ui/widgets/mainform/mainwindow.go`, add `FocusedPanel()` returning the focused `*panel.Model`
- [x] 3.2 In `internal/ui/ui.go`, route the Menu action to the focused panel's `CycleSort` instead of a mock

## 4. Tests

- [x] 4.1 Panel test: the cycle order (Name↑ → Name↓ → Size↑ → … → wrap) and per-panel independence
- [x] 4.2 Panel test: directories above files, ordering within each group, and the `..` row first
- [x] 4.3 Panel test: numeric size ordering (208 before 3.9K)
- [x] 4.4 Panel test: the active column's header shows the direction marker and the selection survives a re-sort
- [x] 4.5 Mainform or ui test: the Menu action cycles the focused panel's sort and shows no mock
- [x] 4.6 Connector test: the Size value is an `app.Size` that displays human-readable text and compares numerically
- [x] 4.7 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 4.8 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 4.9 pty smoke (`make run`): F2 cycles the sort with directories first and a visible header marker, and F10 quits
