## 1. Data contract

- [x] 1.1 In `internal/app/connector.go`, add `Hidden bool` and `Width int` to `Attribute`
- [x] 1.2 In `internal/ui/widgets/panel/panel.go`, skip hidden attributes when building columns (`columnsFor`) and rows (`rowFromAttrList`), use a header attribute's `Width` (falling back to the default), keep hidden data available to `entryIsDir`, and update `setDefaultColumns` and `parentRow` (Name `..`, Size `<DIR>`)

## 2. Filesystem connector

- [x] 2.1 In `internal/connectors/filesystem/connector.go`, replace the header with Name/Size/Date/Time (compact widths) plus a hidden IsDir, and build entry rows with a human-readable size, `<DIR>` for directories, and the modification date and time
- [x] 2.2 Add helpers for the human-readable size and the modification date/time, using `entry.Info()`

## 3. Tests

- [x] 3.1 Filesystem connector test: the header declares Name, Size, Date, Time; a file row has a human-readable size and `YYYY-MM-DD` / `HH:MM:SS`; a directory row shows `<DIR>`; the IsDir attribute is hidden
- [x] 3.2 Panel test: hidden attributes produce no column and no cell; a declared header width is used and the default applies when none is declared; directory detection still works with a hidden IsDir
- [x] 3.3 Panel test: the `..` row shows Name `..` and Size `<DIR>`
- [x] 3.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.6 pty smoke (`make run`): the panel shows Name, Size, Date, Time with directories as `<DIR>`, Enter still descends, and F10 quits
