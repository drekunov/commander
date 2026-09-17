## 1. Theme

- [x] 1.1 In `internal/config/styles.css`, add a `.menu-background` class with `color: #C0C0C0`, `background-color: #000000`, and `font-weight: bold`
- [x] 1.2 In `internal/config/config.go`, add a `MenuBackgroundStyle` field to `Styles` and load it from `.menu-background` in `LoadStyles`

## 2. Buttonbar

- [x] 2.1 In `internal/ui/widgets/buttonbar/buttonbar.go`, add a bar-background base (`MenuBackgroundStyle.Inherit(ButtonBarStyle)`) and use it for the padding between labels, as the base in `numberStyle()`, as the default base in `styleFor()`, and as the base for the active/pressed overlays

## 3. Tests

- [x] 3.1 Config test: `LoadStyles` maps `.menu-background` from the embedded default stylesheet to `MenuBackgroundStyle` with foreground `#C0C0C0`, background `#000000`, and bold
- [x] 3.2 Buttonbar test: the padding and the default number/name styles inherit the menu-background base, and a focused/pressed label keeps it while its foreground stays distinct
- [x] 3.3 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 3.4 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 3.5 pty smoke (`make run`): the bar background between labels is black, each number is light-gray on black, each name is black on light-gray, and F10 quits
