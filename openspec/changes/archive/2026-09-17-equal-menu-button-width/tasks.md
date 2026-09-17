## 1. Buttonbar layout

- [x] 1.1 In `internal/ui/widgets/buttonbar/buttonbar.go`, replace the per-segment width with the fixed layout: a two-cell number field (right-aligned) plus a six-cell name field, with `buttonWidth = min(8, width/10)` and the leftover columns split into equal gaps between adjacent buttons
- [x] 1.2 Keep the row exactly `width` cells, truncating each field to its allotted cells when the button is narrower
- [x] 1.3 Update `actionIndexAtColumn` to map a column to the button whose span (its cells plus the following gap) contains it

## 2. Tests

- [x] 2.1 Buttonbar test: for widths that fit eight-cell buttons, every button spans the same width and the row width equals the terminal width (including a width whose leftover is not divisible by nine)
- [x] 2.2 Buttonbar test: a narrow width shrinks the buttons equally and truncates labels
- [x] 2.3 Buttonbar test: mouse mapping activates the correct button for a click inside a button and inside a gap
- [x] 2.4 Run `go build ./...`, `go test -race ./...`, and `go vet ./...`
- [x] 2.5 Run `make lint` (formatters rewrite in place; golangci-lint v2 `default: all` must pass)
- [x] 2.6 pty smoke (`make run`): the ten buttons have equal width with equal gaps, names are aligned, the row spans the full width, and F10 quits
