# AGENTS.md

Guidance for automated coding agents working in this repository.

## Build & Dev Commands

```bash
make run            # go run ./cmd/commander/main.go
make lint           # gofumpt -w . && gci write && golangci-lint run
make build          # go build -o ./commander with version ldflags
make modup          # go get -u ./... && go mod tidy
make pvs            # pvs-golang analyze .  (PVS-Studio, optional)

go run ./cmd/wm-example   # window manager demo
go test ./...             # no tests exist; passes vacuously
```

## Gotchas

- **`make lint` rewrites files in place** (`gofumpt -w .`, `gci write . --skip-generated -s standard -s default`), then runs `golangci-lint`. It is not read-only. CI instead runs `golangci-lint` v2.4 via the action without the formatters.
- **`make build` injects version info via `-ldflags`** into the unexported vars in `internal/app/version.go` (`version`, `commit`, `branch`, `buildUnixTimestamp`). A plain `go build`/`go run` reports `undefined` for all of them. Keep those var names/paths stable or the ldflags silently stop matching.
- **Lint config is golangci-lint v2** (`.golangci.yml`) with `default: all` — nearly every linter is on. A handful are disabled (`gocritic`, `revive`, `mnd`, `exhaustruct`, `wsl`, etc.); `wsl_v5` is enabled instead of `wsl`. Do not fight linters that CI will flag.
- **CSS theming loads at init()**: `internal/config/styles.css` is `//go:embed`-ed; a `styles.css` in the process working directory overrides it at runtime (`internal/config/config.go:28`). Overrides are relative to the working directory, not the binary.

## Architecture

```
cmd/commander/main.go     app entrypoint (runs internal/ui)
cmd/wm-example/main.go    window manager demo
internal/
  app/                    App struct, Connector interface, UI interface; version vars (ldflags-injected)
  config/                 CSS->lipgloss theme engine (embed: styles.css, parse: css.go, styles: config.go)
  connectors/filesystem/  local filesystem Connector
  ui/                     Bubble Tea Model (F10 quits at ui.go:42), main form, dialogs
  ui/widgets/wm/          compositing window manager (z-index, drag-move, drag-resize; only focused window gets keys)
  ui/widgets/panel/       custom table widget (no evertras/bubble-table dependency)
  ui/widgets/mainform/    main window container
```

## Key Facts

- **Module**: `github.com/drekunov/gc`, Go 1.25. Deps: bubbletea, bubbles, lipgloss, `charmbracelet/x/ansi`, `vanng822/css` (CSS parsing).
- **`Connector.ReadFile` returns `([]byte, error)`, not `string`** (`internal/app/connector.go:13`).
- **No tests exist** — don't assume a test framework; `go test ./...` is a no-op.
- **CI**: GitHub Actions on PRs to `develop`, `self-hosted` runners (`.github/workflows/go.yml`).
- **Branches**: develop is the default/merge target; work happens on topic branches (e.g. `fix/linters`).
- **F10** quits the app (`internal/ui/ui.go:42`); the wm-example has its own F10 handler.
- **F12** dumps the current screen to `screen-<timestamp>.ansi` in the working directory (`internal/ui/dump.go`, key handler at `internal/ui/ui.go`).
