## 1. Data model

- [x] 1.1 In `internal/app/connector.go`, add `Flex bool` to `Attribute`
- [x] 1.2 In `internal/connectors/filesystem/connector.go`, mark the Name header attribute flexible (Size, Date, and Time stay fixed)

## 2. Panel sizing

- [x] 2.1 In `internal/ui/widgets/panel/panel.go`, record each visible column's spec (title, declared width, flex) and add a helper that builds the table columns from the specs and the panel width, giving the flexible column the width left after the fixed columns and padding (minimum eight cells)
- [x] 2.2 Recompute the columns from `SetData` and `SetWidth`, update `contentWidth`, keep fixed widths, and handle the no-flexible and multiple-flexible cases

## 3. Tests

- [x] 3.1 Panel test: with a flexible column and a wide panel, the flexible column fills the remaining width, fixed columns keep their widths, and the content width equals the panel width
- [x] 3.2 Panel test: a narrow panel clamps the flexible column to eight cells while fixed columns keep their widths
- [x] 3.3 Panel test: when no column is flexible, every column uses its declared width
- [x] 3.4 Filesystem connector test: the Name header attribute is marked flexible and Size, Date, and Time are not
- [x] 3.5 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.6 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.7 pty smoke (`make run`): the Name column fills the panel with Size/Date/Time fixed, the layout follows the terminal width, and F10 quits
