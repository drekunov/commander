## 1. Theme

- [x] 1.1 In `internal/config/styles.css`, add a `.menu-number` class with `color: #C0C0C0`, `background-color: #000000`, and `font-weight: bold`
- [x] 1.2 In `internal/config/config.go`, add a `MenuNumberStyle` field to `Styles` and load it from `.menu-number` in `LoadStyles`

## 2. Buttonbar

- [x] 2.1 In `internal/ui/widgets/buttonbar/buttonbar.go`, split each label into its numeric prefix (`strings.TrimSuffix(def.label, def.name)`) and its name; render the number with `MenuNumberStyle` (inheriting the strip) and the name with `styleFor`, then fill the remainder of the segment with the strip style
- [x] 2.2 Keep each segment's exact width: truncate the number to the segment width, then the name to the remaining width, measuring with `lipgloss.Width` and truncating with `ansi.Truncate`

## 3. Tests

- [x] 3.1 Config test: `LoadStyles` maps `.menu-number` from the embedded default stylesheet to `MenuNumberStyle` with foreground `#C0C0C0`, background `#000000`, and bold
- [x] 3.2 Buttonbar test: for a default button, the numeric prefix uses the menu-number style and the name uses the menu-label style
- [x] 3.3 Buttonbar test: for a focused and a pressed button, the numeric prefix keeps the menu-number style while the name uses the active or pressed style
- [x] 3.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.6 pty smoke (`make run`): each bottom-menu number is black-on-light-gray and each name is light-gray-on-black on the blue strip, the row spans the full width, and F10 quits
