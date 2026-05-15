# AGENTS.md

This file provides guidance to automated coding agents working in this repository.

## Build & Dev Commands

```bash
go run ./cmd/commander          # run the main TUI app
go run ./cmd/wm-example         # run the window manager demo
go test ./...                   # run all tests (none yet)
go test -run TestName ./pkg     # run a single test

make lint                       # gofumpt + gci + golangci-lint (v2 config)
make modup                      # go get -u + go mod tidy
make build                      # go build -o ./commander ./cmd/commander/main.go
```

## Lint

Uses `golangci-lint` v2 (`.golangci.yml`). Formatters are `gofumpt` and `gci` (run as separate steps before golangci-lint). The CI workflow also runs `golangci-lint` via the official action.

## Architecture

```
cmd/commander/main.go    -> app entrypoint
cmd/wm-example/main.go   -> window manager demo
internal/
  app/                   -> App struct, Connector interface, UI interface
  config/                -> CSS-based theme engine (embed:styles.css, parse:css.go, styles:config.go)
  connectors/filesystem/ -> local filesystem Connector
  ui/                    -> Bubble Tea Model, dialogs, panels
  ui/widgets/wm/         -> compositing window manager (z-index, drag-move, drag-resize)
  ui/widgets/panel/      -> custom table widget (no evertras/bubble-table dependency)
```

## Key Facts

- **Module**: `github.com/drekunov/gc`, Go 1.25
- **TUI framework**: Bubble Tea (charmbracelet), with Lipgloss for styling
- **CSS theming**: `internal/config/styles.css` is embedded at build time. Placing a `styles.css` in the working directory overrides embedded defaults at runtime (`config.go:29`).
- **Connector.ReadFile** returns `([]byte, error)`, *not* `string`.
- **No tests exist** yet — the CI `go test` step passes vacuously.
- **CI** triggers on PRs to `develop`, runs on `self-hosted` runners.
- **F10** quits the app (handled in `ui/ui.go:43`).
- The window manager (`wm/`) supports mouse-driven move/resize and routes keyboard input only to the focused window.
