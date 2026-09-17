## 1. Theme engine

- [x] 1.1 In `internal/config/css.go`, teach `ApplyStyle` to recognize `font-weight: bold` and apply `Bold(true)`
- [x] 1.2 In `internal/config/styles.css`, add a `.cursor` class and a `.menu-label` class, both black text (`#000000`) on `#C0C0C0` with `font-weight: bold`, so the two are equal
- [x] 1.3 In `internal/config/config.go`, add `CursorStyle` and `MenuLabelStyle` fields to `Styles` and load them from `.cursor` and `.menu-label` in `LoadStyles`

## 2. Panel cursor

- [x] 2.1 In `internal/ui/widgets/panel/panel.go`, build the table `Selected` style from the injected `m.styles.CursorStyle` instead of the hardcoded black-on-gray bold style, preserving the blurred-panel no-highlight branch and the focused full-width right padding

## 3. Bottom-menu labels

- [x] 3.1 In `internal/ui/widgets/buttonbar/buttonbar.go`, render each non-focused, non-pressed button's label text with the injected `MenuLabelStyle` and fill the rest of the segment with `ButtonBarStyle`; keep focused/pressed labels on the existing active/pressed overlay
- [x] 3.2 Ensure the label width is measured with `lipgloss.Width` and truncated so the label plus strip fill exactly the segment width

## 4. Tests

- [x] 4.1 Config test: `ApplyStyle` with `font-weight: bold` produces a bold style, and without it does not
- [x] 4.2 Config test: `LoadStyles` maps `.cursor` and `.menu-label` from the embedded default stylesheet to `CursorStyle` and `MenuLabelStyle` with the expected foreground, background, and bold weight
- [x] 4.3 Panel test: the focused cursor style comes from the injected `CursorStyle`, and blurring clears the highlight
- [x] 4.4 Buttonbar test: a default label renders with the menu-label style on the strip, while focused and pressed labels use the active and pressed styles
- [x] 4.5 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 4.6 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 4.7 pty smoke (`make run`): the panel cursor is black-on-gray bold, the bottom-menu labels are black-on-gray on the blue strip, focused/pressed buttons still look distinct, and F10 quits
