## 1. Panel: header shares the table background

- [x] 1.1 In `internal/ui/widgets/panel/panel.go`, set the table styles' `Header` background from the injected `m.styles.TableStyle` background in the render-style builder used by `View()` (the same builder that applies the focus-aware `Selected`), so the header keeps bold weight/padding but shares the table background
- [x] 1.2 Update the `tableStyles` doc comment to state that the header row shares the table background in both focus states

## 2. Tests

- [x] 2.1 Panel test: assert `tableStyles().Header` background equals the injected `.table` background for a focused panel
- [x] 2.2 Panel test: assert the header background is unchanged (still the table background) after the panel is blurred
- [x] 2.3 Panel test: construct a panel with a custom table background and assert the header uses that same background, covering a `styles.css` override
- [x] 2.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 2.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 2.6 pty smoke (`make run`): both panels' column header rows render with the same background as the listing body, focused and unfocused look consistent, and F10 quits
